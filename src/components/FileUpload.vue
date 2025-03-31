<script setup>
import { ref, onMounted } from 'vue'

const file = ref(null)
const metadata = ref({ key: '', value: '' })
const files = ref([])

const log = (message, details = null) => {
  console.log(`[FileUpload] ${message}`)
  if (details) {
    console.log('[FileUpload] Details:', details)
  }
}

const logError = (message, error, requestDetails = null) => {
  console.error(`[FileUpload] ${message}:`, error)
  if (requestDetails) {
    console.error('[FileUpload] Request details:', requestDetails)
  }
}

const fetchFiles = async () => {
  const query = '{ getFiles { id name size contentType createdAt metadata { key value } } }'
  log("Fetching files from the server...", { query })
  
  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Server responded with ${response.status}: ${errorText}`)
    }

    const result = await response.json()
    if (result.errors) {
      throw new Error(result.errors[0].message)
    }
    
    files.value = result.data.getFiles
    log("Files fetched successfully", { filesCount: files.value.length })
  } catch (error) {
    logError("Error fetching files", error, { query })
    alert('Error fetching files')
  }
}

const uploadFile = async () => {
  if (!file.value) {
    log("No file selected for upload")
    return
  }

  // Läs fil som Base64
  const reader = new FileReader()
  reader.onload = async (e) => {
    const base64Data = e.target.result.split(',')[1] // Ta bort data URL prefix

    const fileInput = {
      name: file.value.name,
      size: file.value.size,
      contentType: file.value.type,
      fileData: base64Data, // Skicka med Base64-kodad fildata
    }
    
    // Lägg endast till metadata om både key och value har värden
    if (metadata.key && metadata.value) {
      fileInput.metadata = [{ key: metadata.key, value: metadata.value }]
    }

    const query = `
      mutation($input: FileInput!) {
        saveFile(input: $input) {
          id
          name
        }
      }
    `

    const variables = {
      input: fileInput
    }

    log(`Uploading file: ${file.value.name}`, { query, variables })

    try {
      const response = await fetch('http://localhost:8080/query', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ query, variables }),
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Server responded with ${response.status}: ${errorText}`)
      }

      const result = await response.json()
      if (result.errors) {
        throw new Error(result.errors[0].message)
      }

      log("File uploaded successfully", { 
        fileId: result.data.saveFile.id,
        fileName: result.data.saveFile.name 
      })
      alert('File uploaded successfully!')
      fetchFiles()
    } catch (error) {
      logError("Error uploading file", error, { query, variables })
      alert('Error uploading file: ' + error.message)
    }
  }

  reader.onerror = (error) => {
    logError("Error reading file", error)
    alert('Error reading file')
  }

  reader.readAsDataURL(file.value)
}

const downloadFile = async (fileId) => {
  const query = `
    query DownloadFile($id: ID!) {
      downloadFile(id: $id) {
        id
        name
        size
        contentType
        createdAt
        fileData
        metadata {
          key
          value
        }
      }
    }
  `

  const variables = { id: fileId }
  log(`Downloading file with ID: ${fileId}`, { query, variables })

  try {
    const response = await fetch('http://localhost:8080/query', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query, variables }),
    })

    if (!response.ok) {
      const errorText = await response.text()
      throw new Error(`Server responded with ${response.status}: ${errorText}`)
    }

    const result = await response.json()
    if (result.errors) {
      throw new Error(result.errors[0].message)
    }
    
    const file = result.data.downloadFile
    if (!file) {
      throw new Error('File not found')
    }

    // Konvertera base64-data till blob och ladda ner filen
    if (file.fileData) {
      const binaryData = atob(file.fileData)
      const bytes = new Uint8Array(binaryData.length)
      for (let i = 0; i < binaryData.length; i++) {
        bytes[i] = binaryData.charCodeAt(i)
      }
      
      const blob = new Blob([bytes], { type: file.contentType })
      const link = document.createElement('a')
      link.href = URL.createObjectURL(blob)
      link.download = file.name
      link.click()

      log(`File ${file.name} downloaded successfully`, { 
        fileId: file.id,
        fileName: file.name,
        fileSize: file.size 
      })
    } else {
      throw new Error('File data is missing')
    }
  } catch (error) {
    logError("Error downloading file", error, { query, variables })
    alert('Error downloading file: ' + error.message)
  }
}

onMounted(() => {
  fetchFiles()
})

const onFileChange = (e) => {
  const selectedFile = e.target.files[0]
  if (selectedFile) {
    file.value = selectedFile
    log(`File selected: ${file.value.name}`)
  } else {
    log("No file selected")
  }
}
</script>

<template>
  <div>
    <h1>Upload a Document</h1>
    <input type="file" @change="onFileChange" />
    <div>
      <label>Metadata Key:</label>
      <input type="text" v-model="metadata.key" />
    </div>
    <div>
      <label>Metadata Value:</label>
      <input type="text" v-model="metadata.value" />
    </div>
    <button @click="uploadFile">Upload</button>

    <h2>Uploaded Files</h2>
    <table class="file-table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Size (bytes)</th>
          <th>Content Type</th>
          <th>Created At</th>
          <th>Metadata</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in files" :key="item.id">
          <td>{{ item.name }}</td>
          <td>{{ item.size }}</td>
          <td>{{ item.contentType }}</td>
          <td>{{ item.createdAt }}</td>
          <td>
            <ul>
              <li v-for="meta in item.metadata" :key="meta.key">
                {{ meta.key }}: {{ meta.value }}
              </li>
            </ul>
          </td>
          <td>
            <button @click="downloadFile(item.id)">Download</button>
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