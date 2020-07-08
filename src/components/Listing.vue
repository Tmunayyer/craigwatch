<template>
  <div id='listingItem'>
  
  <h2 style="text-decoration: bold" >{{ listing.name }}</h2>
  <h4>posted on: {{ listing.date }}</h4>
  <p >{{ listing.content }}</p>
  <h4>posting id: {{ listing.id }} </h4>
 
  </div>
</template>

<script>
  
  export default {
    name: 'Listing',
    props: ['listing', 'index']
  }
</script>
