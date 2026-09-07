#include <cstdio>

#include <webgpu/webgpu.h>

namespace {

constexpr uint32_t kWidth = 900;
constexpr uint32_t kHeight = 600;

WGPUInstance instance;
WGPUSurface surface;
WGPUAdapter selected_adapter;
WGPUDevice device;
WGPUQueue queue;
WGPURenderPipeline pipeline;
WGPUTextureFormat surface_format;

WGPUStringView string_view(const char* text) {
    return {text, WGPU_STRLEN};
}

void draw() {
    WGPUSurfaceTexture surface_texture = WGPU_SURFACE_TEXTURE_INIT;
    wgpuSurfaceGetCurrentTexture(surface, &surface_texture);
    if (surface_texture.status != WGPUSurfaceGetCurrentTextureStatus_SuccessOptimal &&
        surface_texture.status != WGPUSurfaceGetCurrentTextureStatus_SuccessSuboptimal) {
        std::fprintf(stderr, "Unable to acquire a WebGPU surface texture.\n");
        return;
    }

    WGPUTextureView view = wgpuTextureCreateView(surface_texture.texture, nullptr);
    WGPURenderPassColorAttachment color_attachment = WGPU_RENDER_PASS_COLOR_ATTACHMENT_INIT;
    color_attachment.view = view;
    color_attachment.loadOp = WGPULoadOp_Clear;
    color_attachment.storeOp = WGPUStoreOp_Store;
    color_attachment.clearValue = {0.04, 0.06, 0.11, 1.0};

    WGPURenderPassDescriptor pass_descriptor = WGPU_RENDER_PASS_DESCRIPTOR_INIT;
    pass_descriptor.colorAttachmentCount = 1;
    pass_descriptor.colorAttachments = &color_attachment;

    WGPUCommandEncoder encoder = wgpuDeviceCreateCommandEncoder(device, nullptr);
    WGPURenderPassEncoder pass = wgpuCommandEncoderBeginRenderPass(encoder, &pass_descriptor);
    wgpuRenderPassEncoderSetPipeline(pass, pipeline);
    wgpuRenderPassEncoderDraw(pass, 3, 1, 0, 0);
    wgpuRenderPassEncoderEnd(pass);

    WGPUCommandBuffer commands = wgpuCommandEncoderFinish(encoder, nullptr);
    wgpuQueueSubmit(queue, 1, &commands);

    wgpuCommandBufferRelease(commands);
    wgpuRenderPassEncoderRelease(pass);
    wgpuCommandEncoderRelease(encoder);
    wgpuTextureViewRelease(view);
    wgpuTextureRelease(surface_texture.texture);
}

void create_pipeline() {
    constexpr char shader_source[] = R"(
        @vertex fn vs(@builtin(vertex_index) index: u32) -> @builtin(position) vec4f {
            let positions = array<vec2f, 3>(vec2f(0.0, 0.72), vec2f(-0.72, -0.62), vec2f(0.72, -0.62));
            return vec4f(positions[index], 0.0, 1.0);
        }

        @fragment fn fs(@builtin(position) position: vec4f) -> @location(0) vec4f {
            return vec4f(position.x / 900.0, 0.35, position.y / 600.0, 1.0);
        }
    )";

    WGPUShaderSourceWGSL wgsl = WGPU_SHADER_SOURCE_WGSL_INIT;
    wgsl.code = string_view(shader_source);
    WGPUShaderModuleDescriptor shader_descriptor = WGPU_SHADER_MODULE_DESCRIPTOR_INIT;
    shader_descriptor.nextInChain = &wgsl.chain;
    WGPUShaderModule shader = wgpuDeviceCreateShaderModule(device, &shader_descriptor);

    WGPUColorTargetState color_target = WGPU_COLOR_TARGET_STATE_INIT;
    color_target.format = surface_format;

    WGPUVertexState vertex = WGPU_VERTEX_STATE_INIT;
    vertex.module = shader;
    vertex.entryPoint = string_view("vs");
    WGPUFragmentState fragment = WGPU_FRAGMENT_STATE_INIT;
    fragment.module = shader;
    fragment.entryPoint = string_view("fs");
    fragment.targetCount = 1;
    fragment.targets = &color_target;

    WGPUPrimitiveState primitive = WGPU_PRIMITIVE_STATE_INIT;
    primitive.topology = WGPUPrimitiveTopology_TriangleList;
    WGPUMultisampleState multisample = WGPU_MULTISAMPLE_STATE_INIT;
    multisample.count = 1;

    WGPURenderPipelineDescriptor descriptor = WGPU_RENDER_PIPELINE_DESCRIPTOR_INIT;
    descriptor.vertex = vertex;
    descriptor.fragment = &fragment;
    descriptor.primitive = primitive;
    descriptor.multisample = multisample;
    pipeline = wgpuDeviceCreateRenderPipeline(device, &descriptor);
    wgpuShaderModuleRelease(shader);
}

void on_device(WGPURequestDeviceStatus status, WGPUDevice requested_device, WGPUStringView, void*, void*) {
    if (status != WGPURequestDeviceStatus_Success) {
        std::fprintf(stderr, "Could not request a WebGPU device.\n");
        return;
    }
    device = requested_device;
    queue = wgpuDeviceGetQueue(device);

    WGPUSurfaceCapabilities capabilities = WGPU_SURFACE_CAPABILITIES_INIT;
    wgpuSurfaceGetCapabilities(surface, selected_adapter, &capabilities);
    surface_format = capabilities.formats[0];
    wgpuSurfaceCapabilitiesFreeMembers(capabilities);

    WGPUSurfaceConfiguration config = WGPU_SURFACE_CONFIGURATION_INIT;
    config.device = device;
    config.format = surface_format;
    config.width = kWidth;
    config.height = kHeight;
    config.presentMode = WGPUPresentMode_Fifo;
    wgpuSurfaceConfigure(surface, &config);
    wgpuAdapterRelease(selected_adapter);
    selected_adapter = nullptr;

    create_pipeline();
    draw();
}

void on_adapter(WGPURequestAdapterStatus status, WGPUAdapter adapter, WGPUStringView, void*, void*) {
    if (status != WGPURequestAdapterStatus_Success) {
        std::fprintf(stderr, "Could not request a WebGPU adapter. Is WebGPU enabled in this browser?\n");
        return;
    }

    selected_adapter = adapter;
    WGPURequestDeviceCallbackInfo callback = WGPU_REQUEST_DEVICE_CALLBACK_INFO_INIT;
    callback.mode = WGPUCallbackMode_AllowSpontaneous;
    callback.callback = on_device;
    wgpuAdapterRequestDevice(adapter, nullptr, callback);
}

}  // namespace

int main() {
    instance = wgpuCreateInstance(nullptr);

    WGPUEmscriptenSurfaceSourceCanvasHTMLSelector canvas_source =
        WGPU_EMSCRIPTEN_SURFACE_SOURCE_CANVAS_HTML_SELECTOR_INIT;
    canvas_source.selector = string_view("#canvas");
    WGPUSurfaceDescriptor surface_descriptor = WGPU_SURFACE_DESCRIPTOR_INIT;
    surface_descriptor.nextInChain = &canvas_source.chain;
    surface = wgpuInstanceCreateSurface(instance, &surface_descriptor);

    WGPURequestAdapterOptions options = WGPU_REQUEST_ADAPTER_OPTIONS_INIT;
    options.compatibleSurface = surface;
    WGPURequestAdapterCallbackInfo callback = WGPU_REQUEST_ADAPTER_CALLBACK_INFO_INIT;
    callback.mode = WGPUCallbackMode_AllowSpontaneous;
    callback.callback = on_adapter;
    wgpuInstanceRequestAdapter(instance, &options, callback);
    return 0;
}
