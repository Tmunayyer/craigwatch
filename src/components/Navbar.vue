<style scoped>
.navbar {
  background-color: #f3f2f2;
  border-bottom: 1px solid #a7a7a7;
  display: flex;
  flex-direction: row;
  align-content: flex-start;
  padding-top: 5px;
  padding-bottom: 5px;
  width: 100%;
}

ol {
  flex-direction: row;
  display: flex;
  align-items: center;
}

.nav-item {
  padding: 5px 5px 5px 5px;
}

.nav-link {
  text-decoration: none;
}

.nav-logo {
  width: 4em;
  margin-left: 1em;
  margin-right: 1em;
  padding-left: 1em;
}

#spinner-container {
  height: 100%;
}
</style>

<template>
  <div class="navbar">
    <ol class="links">
      <router-link class="nav-link" to="/">
        <img class="nav-logo" src="craigswatchLogo3.png" />
      </router-link>

      <li class="nav-item">
        <router-link class="nav-link" to="/">search</router-link>
      </li>
    </ol>

    <Spinner />
  </div>
</template>

<script>
import Spinner from "./Spinner.vue";

export default {
  name: "Navbar",
  components: {
    Spinner
  }
};
</script>