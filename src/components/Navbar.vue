<template>
  <div>
    <ol class="menu">
      <li class="menuItem">
        <router-link to="/search">search</router-link>
      </li>
      <li class="menuItem">
        <router-link to="/monitor">monitor</router-link>
      </li>
    </ol>
    <img src="craigswatchLogo3.png" />
  </div>
</template>

<script>
export default {
  name: "Navbar"
};
</script>