<script setup>
import { ref, onMounted } from 'vue'
import { RouterView } from 'vue-router'

const isDarkTheme = ref(true) // Mörkt tema som standard

// Funktion för att växla mellan ljust och mörkt tema
const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value
  document.documentElement.setAttribute('data-theme', isDarkTheme.value ? 'dark' : 'light')
}

// Sätt temat när komponenten laddas
onMounted(() => {
  document.documentElement.setAttribute('data-theme', isDarkTheme.value ? 'dark' : 'light')
})
</script>

<template>
  <div class="app-container">
    <header>
      <h1>e-Arkive</h1>
      <nav>
        <ul>
          <li><router-link to="/upload">Upload Document</router-link></li>
          <li><router-link to="/list">Document List</router-link></li>
        </ul>
      </nav>
      <button class="theme-toggle" @click="toggleTheme">
        {{ isDarkTheme ? '🌞 Light Theme' : '🌙 Dark Theme' }}
      </button>
    </header>
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style>
:root {
  --background-color: #ffffff;
  --text-color: #000000;
  --table-border-color: #dddddd;
  --button-bg: #4CAF50;
  --button-hover-bg: #45a049;
  --button-text: white;
  --header-bg: #f8f9fa;
  --border-color: #e1e4e8;
}

[data-theme='dark'] {
  --background-color: #121212;
  --text-color: #e0e0e0;
  --table-border-color: #333333;
  --button-bg: #2c8630;
  --button-hover-bg: #236b25;
  --button-text: #e0e0e0;
  --header-bg: #1e1e1e;
  --border-color: #333333;
}

body {
  margin: 0;
  padding: 0;
  font-family: Arial, sans-serif;
  background-color: var(--background-color);
  color: var(--text-color);
  transition: background-color 0.3s, color 0.3s;
}

.app-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

header {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: var(--header-bg);
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  border: 1px solid var(--border-color);
}

header h1 {
  margin: 0 0 20px 0;
}

nav ul {
  list-style: none;
  padding: 0;
  display: flex;
  gap: 1rem;
  margin-bottom: 20px;
}

nav ul li {
  display: inline;
}

nav ul li a {
  text-decoration: none;
  color: var(--text-color);
  padding: 8px 16px;
  border-radius: 4px;
  border: 1px solid var(--border-color);
  transition: all 0.3s;
}

nav ul li a:hover, nav ul li a.router-link-active {
  background-color: var(--button-bg);
  color: var(--button-text);
}

.theme-toggle {
  padding: 8px 16px;
  background-color: var(--button-bg);
  color: var(--button-text);
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.theme-toggle:hover {
  background-color: var(--button-hover-bg);
}

main {
  padding: 20px;
  background-color: var(--background-color);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

button {
  background-color: var(--button-bg);
  color: var(--button-text);
  border: none;
  padding: 8px 16px;
  margin: 5px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: var(--button-hover-bg);
}

input[type="file"], input[type="text"] {
  padding: 8px;
  margin: 5px 0;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--background-color);
  color: var(--text-color);
  width: 100%;
  max-width: 300px;
}

.file-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  border: 1px solid var(--table-border-color);
}

.file-table th, .file-table td {
  border: 1px solid var(--table-border-color);
  padding: 12px 8px;
  text-align: left;
}

.file-table th {
  background-color: var(--header-bg);
  font-weight: bold;
}

.file-table tr:nth-child(even) {
  background-color: rgba(128, 128, 128, 0.05);
}

.file-table tr:hover {
  background-color: rgba(128, 128, 128, 0.1);
}
</style>
