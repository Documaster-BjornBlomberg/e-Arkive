<script setup>
import { createRouter, createWebHistory, RouterView } from 'vue-router';
// Updated import paths:
import UploadPage from './components/pages/UploadPage.vue';
import ListPage from './components/pages/ListPage.vue';

const routes = [
  { path: '/upload', component: UploadPage },
  { path: '/list', component: ListPage },
  { path: '/', redirect: '/list' },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});
</script>

<template>
  <RouterView />
</template>

<style scoped>
/* Add any necessary styles here */
</style>
