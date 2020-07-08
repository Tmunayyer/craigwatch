<template>
  <div id="app">
  <ol class='menu'>
    <li class='menuItem'><router-link to="/search">search</router-link></li>
    <li class='menuItem'><router-link to="/about">about</router-link></li>
    <li class='menuItem'><router-link to="/monitor">monitor</router-link></li>
  </ol> 
  <img src="craigswatchLogo3.png">
  <router-view></router-view>
  </div>
</template>

<script>

export default {
  name: 'App'
}

</script>

