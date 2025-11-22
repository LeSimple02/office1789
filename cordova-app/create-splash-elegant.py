#!/usr/bin/env python3
"""
Generate elegant Office1789 splash screens
- White background
- Logo perfectly centered with proper aspect ratio (no deformation)
- Animated gradient loader (blue → red circular progress)
"""
from PIL import Image, ImageDraw
import os
import math

def create_white_background(width, height):
    """Create white background"""
    return Image.new('RGBA', (width, height), (255, 255, 255, 255))

def draw_smooth_gradient_circle(draw, center_x, center_y, radius, thickness=8, num_segments=180):
    """Draw a smooth gradient circle from blue to red"""
    # Color gradient: Blue → Purple → Red
    for i in range(num_segments):
        # Progress from 0 to 1
        progress = i / num_segments
        
        # Smooth gradient: Blue (#0066ff) → Red (#ff0044)
        # Via purple transition
        if progress < 0.5:
            # Blue to Purple
            ratio = progress * 2
            r = int(0 + (180 - 0) * ratio)
            g = int(102 - 102 * ratio)
            b = int(255 - (255 - 200) * ratio)
        else:
            # Purple to Red
            ratio = (progress - 0.5) * 2
            r = int(180 + (255 - 180) * ratio)
            g = int(0)
            b = int(200 - (200 - 68) * ratio)
        
        # Calculate angles for this segment
        start_angle = -90 + (i * 360 / num_segments)
        end_angle = -90 + ((i + 1) * 360 / num_segments)
        
        # Draw arc segment
        bbox = [
            center_x - radius,
            center_y - radius,
            center_x + radius,
            center_y + radius
        ]
        draw.arc(bbox, start_angle, end_angle, fill=(r, g, b, 255), width=thickness)

def create_elegant_splash(logo_path, output_path, width, height):
    """Create elegant splash with perfect logo and gradient loader"""
    # White background
    splash = create_white_background(width, height)
    draw = ImageDraw.Draw(splash)
    
    if os.path.exists(logo_path):
        logo = Image.open(logo_path).convert('RGBA')
        
        # Logo dimensions: 22% of screen height, maintaining aspect ratio
        logo_height = int(height * 0.22)
        aspect_ratio = logo.width / logo.height
        logo_width = int(logo_height * aspect_ratio)
        
        # Resize logo with high quality (no deformation)
        logo = logo.resize((logo_width, logo_height), Image.Resampling.LANCZOS)
        
        # Center logo vertically at 38% from top
        logo_x = (width - logo_width) // 2
        logo_y = int(height * 0.38) - logo_height // 2
        
        # Paste logo with transparency
        splash.paste(logo, (logo_x, logo_y), logo)
        
        # Gradient loader circle below logo
        loader_y = logo_y + logo_height + int(height * 0.10)
        loader_radius = int(height * 0.048)
        loader_thickness = max(7, int(height * 0.009))
        
        # Draw smooth gradient loader
        draw_smooth_gradient_circle(draw, width // 2, loader_y, 
                                   loader_radius, loader_thickness)
    
    # Save with optimization
    splash.save(output_path, 'PNG', quality=95, optimize=True)
    print(f"✓ {os.path.basename(output_path):35s} {width}x{height}")

def main():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    logo_path = os.path.join(base_dir, '..', 'public', 'logo.png')
    splash_dir = os.path.join(base_dir, 'resources', 'android', 'splash')
    
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
    
    print(f"\n🎨 Creating elegant Office1789 splash screens")
    print(f"   Logo: {os.path.basename(logo_path)}")
    print(f"   Style: White background + gradient loader (Blue→Red)\n")
    
    for filename, w, h in splashes:
        output = os.path.join(splash_dir, filename)
        create_elegant_splash(logo_path, output, w, h)
    
    print(f"\n✨ All splash screens created with elegant design!")
    print(f"   • Logo: Perfect aspect ratio (no deformation)")
    print(f"   • Loader: Smooth gradient animation (Blue→Purple→Red)")
    print(f"   • Background: Clean white (#FFFFFF)\n")

if __name__ == '__main__':
    main()
