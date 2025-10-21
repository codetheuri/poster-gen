Poster Generator: Template Creation GuideThis document outlines the rules and best practices for creating new HTML poster templates for the Poster Generator application. Following these guidelines will ensure that new templates are consistent, reliable, and work seamlessly with the Go backend.Core ConceptsSelf-Contained: Every template MUST be a single, self-contained HTML file. All CSS and images (as inline SVGs) must be included within the file. No external CSS or image links are allowed, as they will fail in the server-side rendering environment.A4 Dimensions: All templates are rendered to an A4-sized PDF. The HTML and body tags must be styled to fill the A4 dimensions (210mm x 297mm).Dynamic Data: Templates are populated with dynamic data from the Go backend. The backend passes a map of data to the template. Placeholders in the template use Go's html/template syntax (e.g., {{.business_name}}).Case-Sensitive Keys: The placeholder keys in the HTML template are case-sensitive and MUST match the JSON keys defined in the required_fields for that template in the database (e.g., paybill_number, agent_name).Template Structure ChecklistWhen creating a new template, follow this structure precisely.1. Document Head (<head>)DOCTYPE & Charset: Start with <!DOCTYPE html> and <meta charset="UTF-8">.CSS Reset: Include a basic CSS reset (* { margin: 0; padding: 0; box-sizing: border-box; }).Page Size & Margin Fix: This block is mandatory to ensure a perfect, margin-free A4 PDF.@page {
    size: A4;
    margin: 0;
}
html, body {
    width: 210mm;
    height: 297mm;
    margin: 0;
    padding: 0;
    font-family: 'Inter', Arial, sans-serif; /* Recommended font */
}
2. Document Body (<body>)Main Container: The root element should be a single <div class="container"> that is styled to fill the full height and width of the body. Use Flexbox (display: flex; flex-direction: column;) on this container to structure the main sections.Dynamic Data Placeholders:For single values like a business name or agent name, use {{.key_name}}. Example: <div>{{.agent_name}}</div>.For numbers that need to be split into individual boxes, the backend automatically provides a Split version. If your required_fields key is paybill_number, you must use the key paybill_numberSplit in your template loop.<div class="number-boxes">
    {{range .paybill_numberSplit}}
        <div class="digit-box">{{.}}</div>
    {{end}}
</div>
3. Logos and ImagesUse Inline SVG: All logos and icons MUST be embedded as inline <svg> tags. Do not use <img> tags with src attributes pointing to external or local files.Finding SVGs: A good resource for finding and optimizing SVGs for major brands is a vector logo website.Embedding: Copy the optimized SVG code directly into your HTML.Example: Adding the "Equity Bank" TemplateCreate the File: Create a new file named equity-paybill.html in the templates directory.Build the HTML/CSS: Design the poster following all the rules above.Define Required Fields: In your database (or via Postman), create a new entry in the poster_templates table:name: "Equity Bank Paybill"layout: "equity-paybill.html"required_fields:[
  {"name": "paybill_number", "label": "Paybill Number", "type": "text"},
  {"name": "account_number", "label": "Account Number (Optional)", "type": "text"}
]
