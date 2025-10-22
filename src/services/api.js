// --- Configuration ---
// This is the base URL of your running Go backend.
const API_BASE_URL = 'http://localhost:8081/api';
const FILE_SERVER_BASE_URL = 'http://localhost:8081'; // URL for accessing generated files

/**
 * A helper function to handle API responses and errors consistently.
 * @param {Response} response - The raw response from the fetch call.
 * @returns {Promise<any>} - The JSON data from the response.
 * @throws {Error} - Throws an error if the response is not OK.
 */
async function handleResponse(response) {
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({ message: 'An unknown network error occurred.' }));
    throw new Error(errorData.message || 'The server returned an error.');
  }
  return response.json();
}

/**
 * Fetches the list of available poster templates from the backend.
 * @returns {Promise<Array>} - A promise that resolves to an array of template objects.
 */
export async function fetchTemplates() {
  const response = await fetch(`${API_BASE_URL}/posters/templates`);
  const payload = await handleResponse(response);
  
  // Based on your backend's structure, the templates are nested.
  const templates = payload.listdatapayload?.data;

  if (!Array.isArray(templates)) {
    throw new Error("Template data from the API is not in the expected format.");
  }

  // Parse the JSON strings within each template object for easier use in Vue.
  return templates.map(t => {
    if (t && typeof t.required_fields === 'string' && t.required_fields) {
      t.required_fields = JSON.parse(t.required_fields);
    } else {
      t.required_fields = [];
    }
    if (t && typeof t.customization_data === 'string' && t.customization_data) {
      t.customization_data = JSON.parse(t.customization_data);
    } else {
      t.customization_data = {};
    }
    return t;
  });
}

/**
 * Fetches the library of predefined logos from the backend.
 * @returns {Promise<Array>} - A promise that resolves to an array of logo objects.
 */
export async function fetchLogos() {
  const response = await fetch(`${API_BASE_URL}/logos`);
  const payload = await handleResponse(response);
  const logos = payload.datapayload?.data;
  
  if (!Array.isArray(logos)) {
      throw new Error("Logo data from the API is not in the expected format.");
  }
  return logos;
}

/**
 * Sends a request to the backend to generate a new poster.
 * @param {number} templateId - The ID of the template to use.
 * @param {string} businessName - The user's business name.
 * @param {object} data - The dynamic field data (e.g., { "till_number": "123" }).
 * @param {object} customizationData - The styling data (e.g., { "primary_color": "#ff0000" }).
 * @returns {Promise<{pdf_url: string}>} - A promise that resolves to an object containing the full URL of the generated PDF.
 */
export async function generatePoster(templateId, businessName, data, customizationData) {
  const response = await fetch(`${API_BASE_URL}/posters/generate?template_id=${templateId}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      business_name: businessName,
      data: data,
      customization_data: customizationData,
    }),
  });
  
  const payload = await handleResponse(response);
  const pdfUrl = payload.datapayload?.data?.pdf_url;
  
  if (!pdfUrl) {
    throw new Error('PDF URL not found in the server response.');
  }

  // Prepend the file server base URL to the relative PDF path to create a full, clickable link.
  return { pdf_url: `${FILE_SERVER_BASE_URL}/${pdfUrl}` };
}

