<template>
  <div class="product-display">
    <div class="product-container">
      <div class="product-image">
        <img :src="image" :class="{ 'out-of-stock-img': !inStock }">
      </div>
      <div class="product-info">
<!--        <h2>{{ msg }}</h2>-->
        <h1>{{ product }}</h1>
        <p v-if="inStock">In Stock</p>
        <p v-else>Out of Stock</p>

        <ProductDetails :details="details"></ProductDetails>

        <div
            v-for="(variant, index) in variants"
            :key="variant.id"
            @mouseover="updateVariant(index)"
            class="color-circle"
            v-bind:style="{ backgroundColor: variant.color}"
            :class="{ 'out-of-stock-img': !inStock }"
        ></div>
        <button
            class="button"
            :class="{ disabledButton: !inStock }"
            @click="addToCart"
            :disabled="!inStock">
          Add to Cart
        </button>
        <button
            class="button"
            :class="{ disabledButton: !inCart }"
            @click="removeFromCart"
            :disabled="!inCart">
          Remove From Cart
        </button>
      </div>
    </div>
  </div>

</template>

<script>
import ProductDetails from "@/components/ProductDetails";
export default {
  name: "ProductDisplay",
  components: {ProductDetails},
  props: {
    cart: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      product: 'Socks',
      brand: 'Vue Mastery',
      selectedVariant: 0,
      details: ['50% cotton', '30% wool', '20% polyester'],
      variants: [
        { id: 2234, color: 'green', image: './assets/images/socks_green.jpg', quantity: 50 },
        { id: 2235, color: 'blue', image: './assets/images/socks_blue.jpg', quantity: 3 },
      ]
    }
  },
  methods: {
    addToCart() {
      this.$emit('add-to-cart', this.variants[this.selectedVariant].id)
      this.variants[this.selectedVariant].quantity--
    },
    removeFromCart() {
      this.$emit('remove-from-cart', this.variants[this.selectedVariant].id)
      this.variants[this.selectedVariant].quantity++
    },
    updateVariant(index) {
      this.selectedVariant = index
    }
  },
  computed: {
    image() {
      return this.variants[this.selectedVariant].image
    },
    inStock() {
      return this.variants[this.selectedVariant].quantity
    },
    inCart() {
      for (let i = 0; i < this.cart.length; i++) {
        if (this.variants[this.selectedVariant].id === this.cart[i]) {
          return true
        }
      }
      return false;
    }
  }
}
</script>
