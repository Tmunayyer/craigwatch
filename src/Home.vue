<template>
  <div>
    <SearchForm />
    <br />
    <SearchList />
  </div>
</template>

<script>
import SearchForm from "./components/SearchForm.vue";
import SearchList from "./components/SearchList.vue";

export default {
  name: "Home",
  components: {
    SearchForm,
    SearchList
  }
};
</script>
