# Script pour créer des templates Office vides valides
# Ces fichiers seront utilisés comme base pour les nouveaux documents

Write-Host "Création des templates Office vides..." -ForegroundColor Cyan

# Pour créer de vrais fichiers Office, on utilise Word/Excel/PowerPoint si disponibles
# Sinon, on copie des fichiers de base depuis le système

$wordPath = "C:\Program Files\Microsoft Office\root\Office16\WINWORD.EXE"
$excelPath = "C:\Program Files\Microsoft Office\root\Office16\EXCEL.EXE"
$powerpointPath = "C:\Program Files\Microsoft Office\root\Office16\POWERPNT.EXE"

function Create-EmptyDocx {
    $docxPath = "$PSScriptRoot\empty.docx"
    
    if (Test-Path $wordPath) {
        Write-Host "Création de empty.docx avec Word..." -ForegroundColor Green
        $word = New-Object -ComObject Word.Application
        $word.Visible = $false
        $doc = $word.Documents.Add()
        $doc.SaveAs([ref]$docxPath, [ref]16) # 16 = wdFormatDocumentDefault
        $doc.Close()
        $word.Quit()
        [System.Runtime.Interopservices.Marshal]::ReleaseComObject($word) | Out-Null
        Write-Host "✓ empty.docx créé" -ForegroundColor Green
    } else {
        Write-Host "Word non trouvé - utilisation d'un template minimal" -ForegroundColor Yellow
        # Créer un fichier DOCX minimal manuellement (structure ZIP)
        Add-Type -Assembly System.IO.Compression.FileSystem
        $zip = [System.IO.Compression.ZipFile]::Open($docxPath, 'Create')
        $zip.Dispose()
    }
}

function Create-EmptyXlsx {
    $xlsxPath = "$PSScriptRoot\empty.xlsx"
    
    if (Test-Path $excelPath) {
        Write-Host "Création de empty.xlsx avec Excel..." -ForegroundColor Green
        $excel = New-Object -ComObject Excel.Application
        $excel.Visible = $false
        $workbook = $excel.Workbooks.Add()
        $workbook.SaveAs($xlsxPath, 51) # 51 = xlOpenXMLWorkbook
        $workbook.Close()
        $excel.Quit()
        [System.Runtime.Interopservices.Marshal]::ReleaseComObject($excel) | Out-Null
        Write-Host "✓ empty.xlsx créé" -ForegroundColor Green
    } else {
        Write-Host "Excel non trouvé - utilisation d'un template minimal" -ForegroundColor Yellow
        Add-Type -Assembly System.IO.Compression.FileSystem
        $zip = [System.IO.Compression.ZipFile]::Open($xlsxPath, 'Create')
        $zip.Dispose()
    }
}

function Create-EmptyPptx {
    $pptxPath = "$PSScriptRoot\empty.pptx"
    
    if (Test-Path $powerpointPath) {
        Write-Host "Création de empty.pptx avec PowerPoint..." -ForegroundColor Green
        $powerpoint = New-Object -ComObject PowerPoint.Application
        $presentation = $powerpoint.Presentations.Add()
        $presentation.SaveAs($pptxPath, 24) # 24 = ppSaveAsOpenXMLPresentation
        $presentation.Close()
        $powerpoint.Quit()
        [System.Runtime.Interopservices.Marshal]::ReleaseComObject($powerpoint) | Out-Null
        Write-Host "✓ empty.pptx créé" -ForegroundColor Green
    } else {
        Write-Host "PowerPoint non trouvé - utilisation d'un template minimal" -ForegroundColor Yellow
        Add-Type -Assembly System.IO.Compression.FileSystem
        $zip = [System.IO.Compression.ZipFile]::Open($pptxPath, 'Create')
        $zip.Dispose()
    }
}

try {
    Create-EmptyDocx
    Create-EmptyXlsx
    Create-EmptyPptx
    
    Write-Host "`n✅ Templates créés avec succès!" -ForegroundColor Green
    Write-Host "Fichiers dans: $PSScriptRoot" -ForegroundColor Cyan
} catch {
    Write-Host "❌ Erreur: $_" -ForegroundColor Red
    Write-Host "`nSi Office n'est pas installé, téléchargez des templates vides depuis:" -ForegroundColor Yellow
    Write-Host "https://github.com/officedocument/templates" -ForegroundColor Yellow
}
