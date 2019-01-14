<template>
  <div class="container" style="min-width:480px">
    <router-link to="/">Home</router-link>
    <router-link to="/about">About</router-link>
    <router-link to="/dynauth">Authenticate</router-link>
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  data () {
    return {}
  }
}
</script>
<style>
  @import '../assets/css/resume.css';
</style>
