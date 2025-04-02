<script setup>
import { ref, onMounted } from 'vue'

const files = ref([])

const fetchFiles = async () => {
  const query = `
    query {
      getFiles {
        id
        name
        size
        contentType
        createdAt
        metadata {
          key
          value
        }
      }
    }
  `

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
    })

    const result = await response.json()
    if (result.errors) {
      throw new Error(result.errors[0].message)
    }

    files.value = result.data.getFiles
  } catch (error) {
    alert('Error fetching files: ' + error.message)
  }
}

onMounted(fetchFiles)
</script>

<template>
  <div>
    <h1>Document List</h1>
    <table class="file-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Size</th>
          <th>Content Type</th>
          <th>Created At</th>
          <th>Metadata</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="file in files" :key="file.id">
          <td>{{ file.name }}</td>
          <td>{{ file.size }}</td>
          <td>{{ file.contentType }}</td>
          <td>{{ file.createdAt }}</td>
          <td>
            <ul>
              <li v-for="meta in file.metadata" :key="meta.key">
                {{ meta.key }}: {{ meta.value }}
              </li>
            </ul>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.file-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.file-table th, .file-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.file-table th {
  background-color: #f4f4f4;
  font-weight: bold;
}

.file-table tr:nth-child(even) {
  background-color: #f9f9f9;
}

.file-table tr:hover {
  background-color: #f1f1f1;
}
</style>