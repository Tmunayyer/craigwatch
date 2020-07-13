<style module>
#app {
  padding: unset;
  /* min-height: 100%;
  margin: 0px;
  align-items: center;
  position: absolute;
  display: flex;
  flex-direction: column;
  width: 100%; */
}
</style>

<template>
  <div id="app">
    <Navbar />
    <router-view></router-view>
  </div>
</template>

<script>
import Navbar from "./components/Navbar.vue";
export default {
  name: "App",
  components: {
    Navbar
  }
};
</script>

