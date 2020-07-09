<template>
  <div>
    <Navbar />
    <SearchForm />
    <SearchList />
  </div>
</template>

<script>
import Navbar from "./components/Navbar.vue";
import SearchForm from "./components/SearchForm.vue";
import SearchList from "./components/SearchList.vue";

export default {
  name: "Home",
  components: {
    Navbar,
    SearchForm,
    SearchList
  }
};
</script>
