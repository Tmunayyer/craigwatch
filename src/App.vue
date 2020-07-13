<style module>
#app {
  height: inherit;
  width: inherit;

  padding: unset;

  color: rgb(26, 25, 25);
  font-family: avenir;
  font-size: 0.95em;
}

#content {
  box-sizing: border-box;

  height: inherit;
  width: inherit;

  display: flex;
  flex-direction: row;
  flex-grow: 1;
  align-content: left;

  padding: 1em 1.5em 1em 1.5em;
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

