<template>
  <div class="mb-8">
    <h2 class="text-xl font-semibold text-slate-700 border-b-2 border-slate-200 pb-2 mb-4">Step 1: Choose a Template</h2>
    <div v-if="templates.length > 0" class="grid grid-cols-2 md:grid-cols-3 gap-4">
      <div v-for="template in templates" :key="template.id" @click="$emit('select', template)"
        class="cursor-pointer rounded-lg p-4 border-2 transition-all duration-300"
        :class="selectedTemplate && selectedTemplate.id === template.id ? 'border-blue-500 bg-blue-50 shadow-md' : 'border-slate-300 hover:border-blue-400 hover:bg-slate-50'">
        <img :src="template.thumbnail || `https://placehold.co/400x300/e2e8f0/64748b?text=${template.name}`" alt="Template Thumbnail" class="rounded-md mb-2 w-full h-24 object-cover">
        <h3 class="font-semibold text-slate-800 text-center">{{ template.name }}</h3>
      </div>
    </div>
    <div v-else class="text-center text-slate-500 p-4 border-2 border-dashed rounded-lg">
      no templates available
    </div>
  </div>
</template>

<script setup>
defineProps({
  templates: Array,
  selectedTemplate: Object,
});
defineEmits(['select']);
</script>
