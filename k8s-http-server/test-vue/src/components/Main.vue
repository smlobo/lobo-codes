<template>
  <div id="app">
    <div class="nav-bar"></div>

    <div class="cart">Cart({{ cart.length }})</div>

    <ProductDisplay :cart="cart" @add-to-cart="addToCart" @remove-from-cart="removeFromCart">
    </ProductDisplay>
  </div>
</template>

<script>

import ProductDisplay from "@/components/ProductDisplay";

export default {
  name: 'MainPage',
  components: {ProductDisplay},
  props: {
    msg: String
  },
  data() {
    return {
      cart:[],
    }
  },
  methods: {
    addToCart(id) {
      this.cart.push(id)
    },
    removeFromCart(id) {
      this.cart.pop(id)
    }
  }
}

</script>
