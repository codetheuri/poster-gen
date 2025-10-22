<template>
  <div v-if="template" class="mb-6 animate-fade-in">
    <h2 class="text-xl font-semibold text-slate-700 border-b-2 border-slate-200 pb-2 mb-4">Step 2: Customize Your Poster</h2>
    <form @submit.prevent="submitForm" class="space-y-4">
      
      <div>
        <label class="block font-medium text-slate-600 mb-1">Business Name</label>
        <input v-model="formData.businessName" type="text" required class="w-full form-input">
      </div>

      
      <div v-for="field in template.required_fields" :key="field.name">
        <label :for="field.name" class="block font-medium text-slate-600 mb-1">{{ field.label }}</label>
        <input :id="field.name" v-model="formData.data[field.name]" :type="field.type" required class="w-full form-input">
      </div>

      
      <div class="border-t pt-4 mt-4 space-y-4">
          <h3 class="text-lg font-medium text-slate-700">Style Options</h3>
          
          <div>
            <label class="block font-medium text-slate-600 mb-1">Brand Logo</label>
            <select v-model="selectedLogoId" @change="updateCustomizationFromSelection" class="w-full form-input appearance-none bg-white pr-8">
              <option :value="null">-- Select a Brand (Optional) --</option>
              <option v-for="logo in logos" :key="logo.id" :value="logo.id">{{ logo.name }}</option>
            </select>
          </div>

          
          <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block font-medium text-slate-600 mb-1">Primary Color</label>
                <input v-model="formData.customizationData.primary_color" type="color" class="w-full h-10 border border-slate-300 rounded-md p-1 cursor-pointer">
              </div>
              <div>
                <label class="block font-medium text-slate-600 mb-1">Text on Primary</label>
                <input v-model="formData.customizationData.text_color_on_primary" type="color" class="w-full h-10 border border-slate-300 rounded-md p-1 cursor-pointer">
              </div>
              <div>
                <label class="block font-medium text-slate-600 mb-1">Secondary Text</label>
                <input v-model="formData.customizationData.secondary_text_color" type="color" class="w-full h-10 border border-slate-300 rounded-md p-1 cursor-pointer">
              </div>
          </div>

          <div>
            <label class="block font-medium text-slate-600 mb-1">Font Family</label>
            <select v-model="formData.customizationData.font_family_name" class="w-full form-input appearance-none bg-white pr-8">
                <option value="Inter">Inter (Default)</option>
                <option value="Arial">Arial</option>
                <option value="Verdana">Verdana</option>
                <option value="Tahoma">Tahoma</option>
                <option value="Montserrat">Montserrat</option> 
                {/* Add more Google Font names if desired */}
            </select>
          </div>

          <div class="grid grid-cols-3 gap-4">
              <div>
                <label class="block font-medium text-slate-600 mb-1 text-sm">Font Size Large</label>
                <input v-model="formData.customizationData.font_size_large" type="text" placeholder="e.g., 36px" class="w-full form-input text-sm">
              </div>
               <div>
                <label class="block font-medium text-slate-600 mb-1 text-sm">Font Size Medium</label>
                <input v-model="formData.customizationData.font_size_medium" type="text" placeholder="e.g., 28px" class="w-full form-input text-sm">
              </div>
               <div>
                <label class="block font-medium text-slate-600 mb-1 text-sm">Font Size XLarge</label>
                <input v-model="formData.customizationData.font_size_xlarge" type="text" placeholder="e.g., 60px" class="w-full form-input text-sm">
              </div>
          </div>
          
      </div>
      

      
      <button type="submit" :disabled="isLoading" class="w-full btn-primary flex items-center justify-center">
        
        <svg v-if="isLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
        <span>{{ isLoading ? 'Generating...' : 'Generate Poster' }}</span>
      </button>
    </form>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue';

const props = defineProps({
  template: Object, // Can be null initially
  logos: Array,
  isLoading: Boolean,
});
const emit = defineEmits(['generate']);

const formData = ref({
  businessName: '',
  data: {},
  customizationData: {} // Will be populated by the watcher
});
const selectedLogoId = ref(null);

// Initialize or reset form when the template prop changes
watch(() => props.template, (newTemplate) => {
  if (newTemplate) {
      // Start with base customization from the template, ensuring it's an object
      const baseCustomization = newTemplate.customization_data || {};
      
      formData.value = {
        businessName: '',
        data: {},
        // Ensure customizationData is always an object, even if base is empty/null
        // Initialize ALL potential customization fields to ensure reactivity
        customizationData: { 
            primary_color: baseCustomization.primary_color || '#A4002D', // Default red if none provided
            text_color_on_primary: baseCustomization.text_color_on_primary || '#FFFFFF',
            secondary_text_color: baseCustomization.secondary_text_color || '#333333',
            font_family_name: baseCustomization.font_family_name || 'Inter',
            font_size_large: baseCustomization.font_size_large || '36px',
            font_size_medium: baseCustomization.font_size_medium || '28px',
            font_size_xlarge: baseCustomization.font_size_xlarge || '60px',
            header_logo_svg: baseCustomization.header_logo_svg || null,
         } 
      };
      selectedLogoId.value = null; // Reset logo selection
  } else {
      // If template becomes null, clear the form
       formData.value = { businessName: '', data: {}, customizationData: {} };
       selectedLogoId.value = null;
  }
}, { immediate: true, deep: true }); // Use deep watch if template object structure might change

// Update customization when a logo is selected from the dropdown
const updateCustomizationFromSelection = () => {
  const selectedLogo = props.logos.find(logo => logo.id === selectedLogoId.value);
  if (selectedLogo) {
    // Overwrite color and logo SVG with the selected brand's defaults
    // Use defaults only if they exist on the logo object
    if(selectedLogo.default_color) {
        formData.value.customizationData.primary_color = selectedLogo.default_color;
    }
    if(selectedLogo.svg_code) {
        formData.value.customizationData.header_logo_svg = selectedLogo.svg_code;
    }

  } else {
    // If "No Logo" or "-- Select --" is chosen, clear the logo SVG
    // Keep the color as it might have been set manually or from template defaults
     if (formData.value.customizationData) {
       formData.value.customizationData.header_logo_svg = null;
     }
     // Optionally reset the color to the template's default if no logo selected?
     // formData.value.customizationData.primary_color = props.template?.customization_data?.primary_color || '#A4002D';
  }
};

// Emit the complete form data (including customizations) when submitted
const submitForm = () => {
  emit('generate', { 
      businessName: formData.value.businessName, 
      data: formData.value.data, 
      // Send the current state of customizationData
      customizationData: formData.value.customizationData 
  });
};
</script>

<style scoped>
.form-input {
  @apply px-4 py-2 border border-slate-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition w-full;
}
.btn-primary {
  @apply bg-green-600 text-white font-bold py-3 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition-all duration-300 disabled:bg-slate-400 disabled:cursor-not-allowed;
}
/* Add a simple fade-in animation */
.animate-fade-in {
    animation: fadeIn 0.5s ease-out forwards;
}
@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>

