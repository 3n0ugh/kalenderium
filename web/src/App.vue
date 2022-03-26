<template>
  <v-app>
    <v-main>
      <CalendarGo />
    </v-main>
  </v-app>
</template>

<script>
import CalendarGo from './components/CalendarGo';

export default {
  name: 'App',

  components: {
    CalendarGo
  },

  data: () => ({
    //
  }),
};
</script>
