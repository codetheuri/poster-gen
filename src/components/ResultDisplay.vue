<template>
  <div v-if="pdfUrl || errorMessage" class="mt-8 text-center">
    
    <div v-if="pdfUrl" class="bg-green-100 border-l-4 border-green-500 text-green-800 p-4 rounded-md">
      <h3 class="font-bold text-lg mb-2">Poster Generated Successfully!</h3>
      <a :href="pdfUrl" target="_blank"
        class="inline-block bg-green-600 text-white font-bold py-2 px-6 rounded-md hover:bg-green-700 transition-transform hover:scale-105">
        Download PDF
      </a>
    </div>
 
    <div v-if="errorMessage" class="bg-red-100 border-l-4 border-red-500 text-red-800 p-4 rounded-md">
      <h3 class="font-bold">An Error Occurred</h3>
      <p>{{ errorMessage }}</p>
    </div>
  </div>
</template>

<script setup>
defineProps({
  pdfUrl: String,
  errorMessage: String,
});
</script>
