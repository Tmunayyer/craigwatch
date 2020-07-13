<style module>
.navbar {
  background-color: #d5d5d5;
  display: flex;
  flex-direction: row;
  align-content: flex-start;
  padding: 5px;
  width: 100%;
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
}
</style>

<template>
  <ol class="navbar">
    <img class="nav-logo" src="craigswatchLogo3.png" />
    <li class="nav-item">
      <router-link class="nav-link" to="/search">search</router-link>
    </li>
    <li class="nav-item">
      <router-link class="nav-link" to="/monitor">monitor</router-link>
    </li>
  </ol>
</template>

<script>
export default {
  name: "Navbar"
};
</script>