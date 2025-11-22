"""
Script pour scanner tous les fichiers Vue et identifier les textes en dur qui devraient être traduits
"""
import re
from pathlib import Path

def scan_vue_file(file_path):
    """Scanne un fichier Vue pour trouver les textes en dur"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Patterns pour détecter les textes en dur (français ou anglais)
    patterns = [
        # Texte entre guillemets simples ou doubles dans le template (pas $t)
        r'(?<![\$t\(])(["\'])([^"\'$]{3,})\1',
        # Placeholder d'input
        r'placeholder=["\']([^"\']{3,})["\']',
        # Title
        r'title=["\']([^"\']{3,})["\']',
        # aria-label
        r'aria-label=["\']([^"\']{3,})["\']',
    ]
    
    found_texts = []
    for pattern in patterns:
        matches = re.finditer(pattern, content)
        for match in matches:
            text = match.group(2) if len(match.groups()) > 1 else match.group(1)
            # Filtrer les textes qui sont probablement du code/CSS
            if not any(skip in text.lower() for skip in ['http', 'www', 'px', 'rem', 'rgb', 'var(', 'url(', '.', '#', 'function', 'return', 'const', 'let', 'import']):
                # Vérifier que ce n'est pas déjà traduit
                if '$t(' not in content[max(0, match.start()-20):match.start()]:
                    if len(text.strip()) > 3 and text.strip() not in ['svg', 'div', 'span', 'button']:
                        found_texts.append((text, match.start()))
    
    return found_texts

def main():
    print("🔍 Scan des fichiers Vue pour détecter les textes en dur non traduits")
    print("=" * 80)
    
    webfront_path = Path(__file__).parent / "webfront2" / "src"
    vue_files = list(webfront_path.rglob("*.vue"))
    
    all_findings = {}
    
    for vue_file in vue_files:
        relative_path = vue_file.relative_to(webfront_path.parent)
        findings = scan_vue_file(vue_file)
        
        if findings:
            all_findings[str(relative_path)] = findings
    
    if all_findings:
        print(f"\n📝 Trouvé des textes en dur dans {len(all_findings)} fichiers:\n")
        
        for file_path, texts in all_findings.items():
            print(f"\n📄 {file_path}")
            unique_texts = list(set([t[0] for t in texts]))[:10]  # Limiter à 10 exemples
            for text in unique_texts:
                print(f"   • \"{text}\"")
    else:
        print("\n✅ Aucun texte en dur détecté (ou tous sont déjà traduits)")
    
    print("\n" + "=" * 80)
    print(f"Total: {len(vue_files)} fichiers Vue scannés")

if __name__ == "__main__":
    main()
