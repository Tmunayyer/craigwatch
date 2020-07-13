<style module>
#app {
  height: inherit;
  width: inherit;

  padding: unset;
}

#content {
  height: inherit;
  width: inherit;

  display: flex;
  flex-direction: row;
  flex-grow: 1;
  align-content: left;

  margin-top: 1em;
  margin-left: 1.5em;
}
</style>

<template>
  <div id="app">
    <Navbar />
    <div id="content">
      <router-view></router-view>
    </div>
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

