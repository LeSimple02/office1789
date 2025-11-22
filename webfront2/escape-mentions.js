import fs from 'fs';

const data = JSON.parse(fs.readFileSync('src/traduction.json', 'utf8'));

// Escape all [ and ] brackets in mentions (Vue I18n interprets them as placeholders)
data.en.mentions = data.en.mentions
  .replace(/\[/g, "{'['}") 
  .replace(/\]/g, "{']'}");

data.fr.mentions = data.fr.mentions
  .replace(/\[/g, "{'['}")
  .replace(/\]/g, "{']'}")
  .replace(/\{\{date\}\}/g, "{'{{'}date{'}}'}"); // Also escape {{date}}

fs.writeFileSync('src/traduction.json', JSON.stringify(data, null, 4), 'utf8');
console.log('✓ Escaped brackets in mentions (EN & FR)');
