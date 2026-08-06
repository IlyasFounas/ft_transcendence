export default function Avatar() {
  return (
    <svg
      width="52"
      height="64"
      viewBox="0 0 52 64"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <defs>
        <clipPath id="avatarClip">
          <path d="
            M41.4854 21.0381
            V42.9609
            L22.5 53.9229
            L3.51465 42.9609
            V21.0381
            L22.5 10.0762
            L41.4854 21.0381
            Z
          "/>
        </clipPath>

        <filter id="glow">
          <feGaussianBlur stdDeviation="5" />
        </filter>
      </defs>

      <image
        href="/avatars/avatar1.png"
        width="52"
        height="64"
        preserveAspectRatio="xMidYMid meet"
        clipPath="url(#avatarClip)"
      />

      <path
        d="
          M41.4854 21.0381
          V42.9609
          L22.5 53.9229
          L3.51465 42.9609
          V21.0381
          L22.5 10.0762
          L41.4854 21.0381
          Z
        "
        stroke="#35DFFF"
      />

    </svg>
  );
}