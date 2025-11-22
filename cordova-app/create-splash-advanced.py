#!/usr/bin/env python3
"""
Generate advanced splash screens for Office1789
- White background
- Centered logo with proper proportions (square padding)
- Dual-spinner loader (blue/red gradient)
"""
from PIL import Image, ImageDraw
import os
import math

def create_white_background(width, height):
    """Create white background"""
    return Image.new('RGBA', (width, height), (255, 255, 255, 255))

def draw_gradient_spinner(draw, center_x, center_y, radius, start_angle, end_angle, color1, color2, thickness=8):
    """Draw a gradient arc segment"""
    steps = 50
    angle_step = (end_angle - start_angle) / steps
    
    for i in range(steps):
        # Interpolate color
        ratio = i / steps
        r = int(color1[0] + (color2[0] - color1[0]) * ratio)
        g = int(color1[1] + (color2[1] - color1[1]) * ratio)
        b = int(color1[2] + (color2[2] - color1[2]) * ratio)
        
        angle = start_angle + i * angle_step
        next_angle = angle + angle_step
        
        # Draw arc segment
        bbox = [
            center_x - radius, 
            center_y - radius,
            center_x + radius,
            center_y + radius
        ]
        draw.arc(bbox, angle, next_angle, fill=(r, g, b, 255), width=thickness)

def create_splash_with_loader(logo_path, output_path, width, height):
    """Create splash with logo and animated-style loader"""
    # White background
    splash = create_white_background(width, height)
    draw = ImageDraw.Draw(splash)
    
    # Load logo
    if os.path.exists(logo_path):
        logo = Image.open(logo_path).convert('RGBA')
        
        # Logo size: 25% of screen height, maintain square aspect
        logo_size = int(height * 0.25)
        
        # Resize logo proportionally within square bounds
        logo.thumbnail((logo_size, logo_size), Image.Resampling.LANCZOS)
        
        # Center logo in upper portion (40% from top)
        logo_x = (width - logo.width) // 2
        logo_y = int(height * 0.40) - logo.height // 2
        
        splash.paste(logo, (logo_x, logo_y), logo)
        
        # Dual spinner below logo
        spinner_y = logo_y + logo.height + int(height * 0.08)
        spinner_radius = int(height * 0.045)
        
        # Blue arc (0° to 180°)
        color_blue_start = (0, 102, 255)    # #0066ff
        color_blue_end = (0, 78, 146)       # #004e92
        draw_gradient_spinner(draw, width // 2, spinner_y, spinner_radius, 
                            -90, 90, color_blue_start, color_blue_end, 
                            thickness=max(6, int(height * 0.008)))
        
        # Red arc (180° to 360°)
        color_red_start = (255, 0, 68)      # #ff0044
        color_red_end = (200, 0, 50)        # #c80032
        draw_gradient_spinner(draw, width // 2, spinner_y, spinner_radius, 
                            90, 270, color_red_start, color_red_end, 
                            thickness=max(6, int(height * 0.008)))
    
    # Save splash
    splash.save(output_path, 'PNG', quality=95, optimize=True)
    print(f"✓ {output_path} ({width}x{height})")

def main():
    base_dir = os.path.dirname(os.path.abspath(__file__))
    logo_path = os.path.join(base_dir, '..', 'public', 'logo.png')
    splash_dir = os.path.join(base_dir, 'resources', 'android', 'splash')
    
    os.makedirs(splash_dir, exist_ok=True)
    
    # Android splash dimensions
    splashes = [
        ('screen-ldpi-portrait.png', 320, 426),
        ('screen-mdpi-portrait.png', 320, 470),
        ('screen-hdpi-portrait.png', 480, 640),
        ('screen-xhdpi-portrait.png', 720, 960),
        ('screen-xxhdpi-portrait.png', 960, 1600),
        ('screen-xxxhdpi-portrait.png', 1280, 1920),
    ]
    
    print(f"Creating Office1789 splash screens (white bg + gradient loader)")
    print(f"Logo: {logo_path}\n")
    
    for filename, w, h in splashes:
        output = os.path.join(splash_dir, filename)
        create_splash_with_loader(logo_path, output, w, h)
    
    print(f"\n✓ All splash screens created!")

if __name__ == '__main__':
    main()
