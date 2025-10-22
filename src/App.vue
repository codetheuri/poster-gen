<template>
  <div class="bg-slate-100 flex items-center justify-center min-h-screen p-4 font-sans">
    <div class="w-full max-w-3xl">
      <div class="bg-white rounded-xl shadow-lg p-6 md:p-8">
        <header class="text-center mb-8">
          <h1 class="text-4xl font-bold text-slate-800">Poster Generator</h1>
          <p class="text-slate-500 mt-2">Create professional posters in seconds.</p>
        </header>

        <TemplateSelector 
          :templates="templates" 
          :selected-template="selectedTemplate"
          @select="handleTemplateSelection"
        />

        <PosterForm 
          v-if="selectedTemplate"
          :template="selectedTemplate"
          :logos="logos"
          :is-loading="isLoading"
          @generate="handleGeneratePoster"
        />

        <ResultDisplay 
          :pdf-url="pdfUrl"
          :error-message="errorMessage"
        />

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import TemplateSelector from './components/TemplateSelector.vue';
import PosterForm from './components/PosterForm.vue';
import ResultDisplay from './components/ResultDisplay.vue';
import { fetchTemplates, fetchLogos, generatePoster as apiGeneratePoster } from './services/api';

const templates = ref([]);
const logos = ref([]);
const selectedTemplate = ref(null);
const isLoading = ref(false);
const pdfUrl = ref('');
const errorMessage = ref('');

onMounted(async () => {
  try {
    templates.value = await fetchTemplates();
    logos.value = await fetchLogos();
  } catch (error) {
    errorMessage.value = error.message;
  }
});

const handleTemplateSelection = (template) => {
  selectedTemplate.value = template;
  pdfUrl.value = '';
  errorMessage.value = '';
};

const handleGeneratePoster = async ({ businessName, data, customizationData }) => {
  isLoading.value = true;
  pdfUrl.value = '';
  errorMessage.value = '';

  try {
    const response = await apiGeneratePoster(selectedTemplate.value.id, businessName, data, customizationData);
    pdfUrl.value = response.pdf_url;
  } catch (error) {
    errorMessage.value = error.message;
  } finally {
    isLoading.value = false;
  }
};
</script>
