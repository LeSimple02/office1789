#!/usr/bin/env python3
"""
Generate custom splash screens for Cordova Android with Office1789 branding
Background: Blue gradient (#000428 -> #004e92)
Logo: Centered with proper scaling
"""
from PIL import Image, ImageDraw
import os

def create_gradient_background(width, height):
    """Create a blue gradient background"""
    img = Image.new('RGB', (width, height))
    draw = ImageDraw.Draw(img)
    
    # Gradient colors
    color1 = (0, 4, 40)   # #000428
    color2 = (0, 78, 146) # #004e92
    
    for y in range(height):
        # Linear interpolation between colors
        ratio = y / height
        r = int(color1[0] + (color2[0] - color1[0]) * ratio)
        g = int(color1[1] + (color2[1] - color1[1]) * ratio)
        b = int(color1[2] + (color2[2] - color1[2]) * ratio)
        draw.line([(0, y), (width, y)], fill=(r, g, b))
    
    return img

def create_splash(logo_path, output_path, width, height):
    """Create splash screen with logo centered on gradient"""
    # Create gradient background
    splash = create_gradient_background(width, height)
    
    # Load and resize logo
    if os.path.exists(logo_path):
        logo = Image.open(logo_path).convert('RGBA')
        
        # Calculate logo size (30% of screen height)
        logo_height = int(height * 0.3)
        aspect_ratio = logo.width / logo.height
        logo_width = int(logo_height * aspect_ratio)
        
        logo = logo.resize((logo_width, logo_height), Image.Resampling.LANCZOS)
        
        # Center logo
        x = (width - logo_width) // 2
        y = (height - logo_height) // 2
        
        # Paste logo with alpha channel
        splash.paste(logo, (x, y), logo)
    
    # Save splash
    splash.save(output_path, 'PNG', quality=95)
    print(f"✓ Created {output_path} ({width}x{height})")

def main():
    # Paths
    base_dir = os.path.dirname(os.path.abspath(__file__))
    logo_path = os.path.join(base_dir, '..', 'public', 'logo.png')
    splash_dir = os.path.join(base_dir, 'resources', 'android', 'splash')
    
    # Ensure splash directory exists
    os.makedirs(splash_dir, exist_ok=True)
    
    # Android splash dimensions (portrait)
    splashes = [
        ('screen-ldpi-portrait.png', 320, 426),
        ('screen-mdpi-portrait.png', 320, 470),
        ('screen-hdpi-portrait.png', 480, 640),
        ('screen-xhdpi-portrait.png', 720, 960),
        ('screen-xxhdpi-portrait.png', 960, 1600),
        ('screen-xxxhdpi-portrait.png', 1280, 1920),
    ]
    
    print(f"Logo source: {logo_path}")
    print(f"Output directory: {splash_dir}\n")
    
    for filename, width, height in splashes:
        output_path = os.path.join(splash_dir, filename)
        create_splash(logo_path, output_path, width, height)
    
    print(f"\n✓ All splash screens created successfully!")
    print(f"\nNext steps:")
    print(f"1. Update config.xml with splash references")
    print(f"2. Run: npm run cordova:android")
    print(f"3. Install APK on device")

if __name__ == '__main__':
    main()
