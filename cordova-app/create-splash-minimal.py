#!/usr/bin/env python3
"""
Generate minimal Office1789 splash screens
- White background
- Logo perfectly centered with proper aspect ratio (no deformation)
- No loader, no extra elements - just clean and simple
"""
from PIL import Image
import os

def create_minimal_splash(logo_path, output_path, width, height):
    """Create minimal splash with only logo on white background"""
    # White background
    splash = Image.new('RGBA', (width, height), (255, 255, 255, 255))
    
    if os.path.exists(logo_path):
        logo = Image.open(logo_path).convert('RGBA')
        
        # Logo dimensions: fit within 35% of screen height AND width, maintaining perfect aspect ratio
        max_logo_height = int(height * 0.35)
        max_logo_width = int(width * 0.55)  # Max 55% of width
        
        # Calculate resize to fit within both constraints
        aspect_ratio = logo.width / logo.height
        
        # Try height-based sizing first
        logo_height = max_logo_height
        logo_width = int(logo_height * aspect_ratio)
        
        # If width exceeds limit, resize based on width instead
        if logo_width > max_logo_width:
            logo_width = max_logo_width
            logo_height = int(logo_width / aspect_ratio)
        
        # Resize logo with highest quality (no deformation)
        logo = logo.resize((logo_width, logo_height), Image.Resampling.LANCZOS)
        
        # Center logo perfectly
        logo_x = (width - logo_width) // 2
        logo_y = (height - logo_height) // 2
        
        # Paste logo with transparency
        splash.paste(logo, (logo_x, logo_y), logo)
    
    # Save
    splash.save(output_path, 'PNG', quality=95, optimize=True)
    print(f"✓ {os.path.basename(output_path):35s} {width}x{height}")

def main():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    logo_path = os.path.join(base_dir, 'public', 'logo.png')  # Changed from '../public'
    splash_dir = os.path.join(base_dir, 'resources', 'android', 'splash')
    
    os.makedirs(splash_dir, exist_ok=True)
    
    splashes = [
        ('screen-ldpi-portrait.png', 320, 426),
        ('screen-mdpi-portrait.png', 320, 470),
        ('screen-hdpi-portrait.png', 480, 640),
        ('screen-xhdpi-portrait.png', 720, 960),
        ('screen-xxhdpi-portrait.png', 960, 1600),
        ('screen-xxxhdpi-portrait.png', 1280, 1920),
    ]
    
    print(f"\n✨ Creating minimal Office1789 splash screens")
    print(f"   Logo: {os.path.basename(logo_path)}")
    print(f"   Style: Logo only, white background\n")
    
    for filename, w, h in splashes:
        output = os.path.join(splash_dir, filename)
        create_minimal_splash(logo_path, output, w, h)
    
    print(f"\n✓ Done! Clean splash with perfect logo proportions.\n")

if __name__ == '__main__':
    main()
