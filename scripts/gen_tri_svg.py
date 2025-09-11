#!/usr/bin/env python3
"""
SVG Triangle Animation Generator (4 Phases, Looping)
Phase 1: build-in (staggered grow/rotate inward)
Phase 2: peel-in  (outer -> inner collapse; hide as each finishes)
Phase 3: peel-out (inner -> outer expand; show as each starts)  [inverse of P2]
Phase 4: unwind   (staggered return to initial)                 [inverse of P1]

All phases loop in sync via a repeating "clock" animate.
"""

def generate_triangle_svg():
    """Generate complete SVG with 4-phase triangle animation (looping)"""
    # Location & Scale
    x_val = 905
    y_val = 700
    overall_scale = 0.7

    # Animation parameters
    num_triangles = 6
    base_delay = 0.35        # initial delay before P1 staggering starts
    base_duration = 0.8      # per-step duration (used for all phases)
    scale_factor = 0.75
    rotation_increment = 12

    # Phase durations (in seconds)
    # P1 staggers over (num_triangles-1) steps after base_delay
    phase1_duration = base_delay + base_duration * (num_triangles - 1)
    phase2_duration = base_duration * (num_triangles - 1)  # one step per i=1..N-1
    phase3_duration = base_duration * (num_triangles - 1)  # inverse order
    phase4_duration = base_duration * (num_triangles - 2) if num_triangles > 2 else 0.0
    # (P4 has steps for i=2..N-1 returning to their previous level)

    # Starts
    phase1_start = 0.0
    phase2_start = phase1_start + phase1_duration
    phase3_start = phase2_start + phase2_duration
    phase4_start = phase3_start + phase3_duration
    total_duration = phase4_start + phase4_duration  # one full cycle

    svg_content = '''{{/* web/templates/tri_anim.html.tmpl */}}
{{define "tri_anim"}}
<g transform="translate(%d %d) scale(%g)">
  <!-- Loop clock: everything else begins relative to this and re-triggers on repeats -->
  <animate id="clock" attributeName="opacity" from="1" to="1"
           begin="0s" dur="%.2fs" repeatCount="indefinite" />
''' % (x_val, y_val, overall_scale, total_duration)

    # Precompute rotation/scale per level i (1..N)
    rot_scale_vals = []
    rot = 0.0
    scale = 1.0
    for _ in range(1, num_triangles + 1):
        rot_scale_vals.append((rot, scale))
        rot += rotation_increment
        scale *= scale_factor

    # Emit triangles with nested rotate/scale groups
    for i in range(1, num_triangles + 1):
        triangle_id = f"n{i}"
        svg_content += f'  <!-- Triangle {i} -->\n'
        svg_content += f'  <g id="rotate_{triangle_id}">\n'
        svg_content += f'    <g id="scale_{triangle_id}">\n'

        # Prevent first-frame flash for i>1 by initializing opacity=0
        opacity_attr = ' opacity="0"' if i > 1 else ''
        svg_content += f'      <use href="#tri" id="{triangle_id}"{opacity_attr}>\n'

        # PHASE 1 (build-in): triangles 2..N appear staggered and move from level i-1 -> i
        if i > 1:
            p1_begin = base_delay + base_duration * (i - 2)
            svg_content += (
                f'        <animate attributeName="opacity" from="0" to="1"\n'
                f'                 begin="clock.begin+{p1_begin:.2f}s; clock.repeatEvent+{p1_begin:.2f}s" dur=".01s" fill="freeze" />\n'
            )

        # PHASE 2 visibility: hide at the end of each P2 step (i=1..N-1)
        if i < num_triangles:
            p2_begin = phase2_start + base_duration * (i - 1)
            p2_end = p2_begin + base_duration
            svg_content += (
                f'        <animate attributeName="opacity" from="1" to="0"\n'
                f'                 begin="clock.begin+{p2_end:.2f}s; clock.repeatEvent+{p2_end:.2f}s" dur=".01s" fill="freeze" />\n'
            )

        # PHASE 3 visibility: re-show as its P3 step begins (inverse order)
        if i < num_triangles:
            k = (num_triangles - 1) - i
            p3_begin = phase3_start + base_duration * k
            svg_content += (
                f'        <animate attributeName="opacity" from="0" to="1"\n'
                f'                 begin="clock.begin+{p3_begin:.2f}s; clock.repeatEvent+{p3_begin:.2f}s" dur=".01s" fill="freeze" />\n'
            )

        # PHASE 4 visibility: return to original (only outermost visible at cycle end)
        if i > 1 and phase4_duration > 0:
            k4 = (num_triangles - 1) - i  # 0-based from i=N-1 down to 2
            p4_end = phase4_start + base_duration * (k4 + 1)
            svg_content += (
                f'        <animate attributeName="opacity" from="1" to="0"\n'
                f'                 begin="clock.begin+{p4_end:.2f}s; clock.repeatEvent+{p4_end:.2f}s" dur=".01s" fill="freeze" />\n'
            )

        svg_content += f'      </use>\n'

        # ---- SCALE GROUP ANIMATIONS ----
        # P1: for i>1, scale from level i-1 -> i
        if i > 1:
            p1_begin = base_delay + base_duration * (i - 2)
            start_scale1 = rot_scale_vals[i - 2][1]
            final_scale1 = rot_scale_vals[i - 1][1]
            svg_content += (
                f'      <animateTransform attributeName="transform" type="scale"\n'
                f'                        begin="clock.begin+{p1_begin:.2f}s; clock.repeatEvent+{p1_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                        from="{start_scale1:.8f}" to="{final_scale1:.8f}" fill="freeze" />\n'
            )

        # P2: for i=1..N-1, scale from level i -> i+1
        if i < num_triangles:
            p2_begin = phase2_start + base_duration * (i - 1)
            start_scale2 = rot_scale_vals[i - 1][1]
            final_scale2 = rot_scale_vals[i][1]
            svg_content += (
                f'      <animateTransform attributeName="transform" type="scale"\n'
                f'                        begin="clock.begin+{p2_begin:.2f}s; clock.repeatEvent+{p2_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                        from="{start_scale2:.8f}" to="{final_scale2:.8f}" fill="freeze" />\n'
            )

        # P3 (inverse of P2): reverse order, scale from level i -> i-1
        if i < num_triangles:
            k = (num_triangles - 1) - i                    # 0..(N-2)
            p3_begin = phase3_start + base_duration * k
            start_scale3 = rot_scale_vals[i][1]
            final_scale3 = rot_scale_vals[i - 1][1]
            svg_content += (
                f'      <animateTransform attributeName="transform" type="scale"\n'
                f'                        begin="clock.begin+{p3_begin:.2f}s; clock.repeatEvent+{p3_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                        from="{start_scale3:.8f}" to="{final_scale3:.8f}" fill="freeze" />\n'
            )

        # P4 (inverse of P1): i=2..N-1 back one level (staggered reverse)
        if 1 < i <= num_triangles and phase4_duration > 0:
            k4 = (num_triangles - 1) - i                   # 0..(N-3)
            p4_begin = phase4_start + base_duration * k4
            start_scale4 = rot_scale_vals[i - 1][1]
            final_scale4 = rot_scale_vals[i - 2][1]
            svg_content += (
                f'      <animateTransform attributeName="transform" type="scale"\n'
                f'                        begin="clock.begin+{p4_begin:.2f}s; clock.repeatEvent+{p4_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                        from="{start_scale4:.8f}" to="{final_scale4:.8f}" fill="freeze" />\n'
            )

        svg_content += f'    </g>\n'

        # ---- ROTATION GROUP ANIMATIONS ----
        # P1 rotate for i>1: level i-1 -> i
        if i > 1:
            p1_begin = base_delay + base_duration * (i - 2)
            start_rot1 = rot_scale_vals[i - 2][0]
            final_rot1 = rot_scale_vals[i - 1][0]
            svg_content += (
                f'  <animateTransform attributeName="transform" type="rotate"\n'
                f'                    begin="clock.begin+{p1_begin:.2f}s; clock.repeatEvent+{p1_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                    from="{start_rot1} 0 0" to="{final_rot1} 0 0" fill="freeze" />\n'
            )

        # P2 rotate for i=1..N-1: level i -> i+1
        if i < num_triangles:
            p2_begin = phase2_start + base_duration * (i - 1)
            start_rot2 = rot_scale_vals[i - 1][0]
            final_rot2 = rot_scale_vals[i][0]
            svg_content += (
                f'  <animateTransform attributeName="transform" type="rotate"\n'
                f'                    begin="clock.begin+{p2_begin:.2f}s; clock.repeatEvent+{p2_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                    from="{start_rot2} 0 0" to="{final_rot2} 0 0" fill="freeze" />\n'
            )

        # P3 (inverse of P2) rotate: reverse order, level i -> i-1
        if i < num_triangles:
            k = (num_triangles - 1) - i
            p3_begin = phase3_start + base_duration * k
            start_rot3 = rot_scale_vals[i][0]
            final_rot3 = rot_scale_vals[i - 1][0]
            svg_content += (
                f'  <animateTransform attributeName="transform" type="rotate"\n'
                f'                    begin="clock.begin+{p3_begin:.2f}s; clock.repeatEvent+{p3_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                    from="{start_rot3} 0 0" to="{final_rot3} 0 0" fill="freeze" />\n'
            )

        # P4 (inverse of P1) rotate: i=2..N-1 back one level (staggered reverse)
        if 1 < i <= num_triangles and phase4_duration > 0:
            k4 = (num_triangles - 1) - i
            p4_begin = phase4_start + base_duration * k4
            start_rot4 = rot_scale_vals[i - 1][0]
            final_rot4 = rot_scale_vals[i - 2][0]
            svg_content += (
                f'  <animateTransform attributeName="transform" type="rotate"\n'
                f'                    begin="clock.begin+{p4_begin:.2f}s; clock.repeatEvent+{p4_begin:.2f}s" dur="{base_duration:.2f}s"\n'
                f'                    from="{start_rot4} 0 0" to="{final_rot4} 0 0" fill="freeze" />\n'
            )

        svg_content += f'  </g>\n\n'

    svg_content += '  </g>\n{{end}}'
    return svg_content


def save_svg(filename="./web/templates/tri_anim.html.tmpl"):
    svg_content = generate_triangle_svg()
    with open(filename, 'w') as f:
        f.write(svg_content)

    # Mirror the timings here for a quick console summary
    num_triangles = 6
    base_delay = 0.35
    base_duration = 0.8
    phase1_duration = base_delay + base_duration * (num_triangles - 1)
    phase2_duration = base_duration * (num_triangles - 1)
    phase3_duration = base_duration * (num_triangles - 1)
    phase4_duration = base_duration * (num_triangles - 2) if num_triangles > 2 else 0.0

    phase1_start = 0.0
    phase2_start = phase1_start + phase1_duration
    phase3_start = phase2_start + phase2_duration
    phase4_start = phase3_start + phase3_duration
    total_duration = phase4_start + phase4_duration

    print(f"SVG animation saved to {filename}")
    print(f"Total animation duration (cycle): {total_duration:.2f} seconds\n")
    print("Animation phases:")
    print(f"Phase 1 ({phase1_start:.2f}s - {phase2_start:.2f}s): build-in (staggered inward)")
    print(f"Phase 2 ({phase2_start:.2f}s - {phase3_start:.2f}s): peel-in (outer→inner, hide each)")
    print(f"Phase 3 ({phase3_start:.2f}s - {phase4_start:.2f}s): peel-out (inner→outer, show each)")
    print(f"Phase 4 ({phase4_start:.2f}s - {total_duration:.2f}s): unwind (stagger back to original)")

if __name__ == "__main__":
    save_svg()

