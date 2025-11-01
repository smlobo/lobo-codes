window.addEventListener("DOMContentLoaded", function () {
    const dark_mode_btn = document.getElementById("dark_mode_btn");
    const light_mode_btn = document.getElementById("light_mode_btn");

    dark_mode_btn.addEventListener('click', function () {
        document.documentElement.setAttribute("data-theme", "dark");
        localStorage.theme = 'dark';
        const mapElement = document.getElementById("mapFrame");
        if (mapElement) {
            mapElement.contentWindow.postMessage("theme-dark", "*");
        }
    });

    light_mode_btn.addEventListener('click', function () {
        document.documentElement.setAttribute("data-theme", "light");
        localStorage.theme = 'light';
        const mapElement = document.getElementById("mapFrame");
        if (mapElement) {
            mapElement.contentWindow.postMessage("theme-light", "*");
        }
    });
});
