import { useId } from "react";

export function SmallButtonSvg({ className }: { className?: string }) {
  const uid = useId();
  const gradientId = `primarySmallGradient-${uid}`;
  const filterId = `primarySmallFilter-${uid}`;

  return (
    <svg width="254" height="98" viewBox="0 0 254 98" fill="none" className={className}>
		<g filter={`url(#${filterId})`}>
		<path d="M20.52 21.096L24.04 17H229.96L233.48 21.096L237 25.192V72.808L233.48 76.904L229.96 81H24.04L20.6642 77.0718L17 72.808V25.192L20.52 21.096Z" fill={`url(#${gradientId})`}/>
		<path d="M229.501 18L236 25.5625V72.4365L229.501 80H24.499L21.4229 76.4199L18 72.4365V25.5625L24.499 18H229.501Z" stroke="#64D0E8" stroke-width="2"/>
		</g>
		<defs>
		<filter id={filterId} x="0" y="0" width="254" height="98" filterUnits="userSpaceOnUse" color-interpolation-filters="sRGB">
		<feFlood flood-opacity="0" result="BackgroundImageFix"/>
		<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
		<feOffset/>
		<feGaussianBlur stdDeviation="8.5"/>
		<feComposite in2="hardAlpha" operator="out"/>
		<feColorMatrix type="matrix" values="0 0 0 0 0.156863 0 0 0 0 0.843137 0 0 0 0 1 0 0 0 0.3 0"/>
		<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_298_338"/>
		<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_298_338" result="shape"/>
		<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
		<feOffset/>
		<feGaussianBlur stdDeviation="11"/>
		<feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1"/>
		<feColorMatrix type="matrix" values="0 0 0 0 0.223529 0 0 0 0 0.85098 0 0 0 0 1 0 0 0 0.25 0"/>
		<feBlend mode="normal" in2="shape" result="effect2_innerShadow_298_338"/>
		</filter>
		<linearGradient id={gradientId} x1="50" y1="118.12" x2="246.107" y2="39.8035" gradientUnits="userSpaceOnUse">
			<stop stopColor="var(--button-gradient-start)"/>
			<stop offset="0.570524" stopColor="var(--button-gradient-middle)"/>
			<stop offset="1" stopColor="var(--button-gradient-end)"/>
		</linearGradient>
		</defs>
    </svg>
  );
}

export function MediumButtonSvg({ className }: { className?: string }) {
  const uid = useId();
  const gradientId = `primaryMediumGradient-${uid}`;
  const filterId = `primaryMediumFilter-${uid}`;

  return (
   <svg width="394" height="124" viewBox="0 0 394 124" fill="none" className={className}>
		<g filter={`url(#${filterId})`}>
		<path d="M22.76 22.76L28.52 17H365.48L371.24 22.76L377 28.52V95.48L371.24 101.24L365.48 107H28.52L22.996 101.476L17 95.48V28.52L22.76 22.76Z" fill={`url(#${gradientId})`}/>
		<path d="M365.065 18L376 28.9346V95.0654L365.065 106H28.9346L18 95.0654V28.9346L28.9346 18H365.065Z" stroke="#64D0E8" stroke-width="2"/>
		</g>
		<defs>
		<filter id={filterId} x="0" y="0" width="394" height="124" filterUnits="userSpaceOnUse" color-interpolation-filters="sRGB">
		<feFlood flood-opacity="0" result="BackgroundImageFix"/>
		<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
		<feOffset/>
		<feGaussianBlur stdDeviation="8.5"/>
		<feComposite in2="hardAlpha" operator="out"/>
		<feColorMatrix type="matrix" values="0 0 0 0 0.156863 0 0 0 0 0.843137 0 0 0 0 1 0 0 0 0.3 0"/>
		<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_298_336"/>
		<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_298_336" result="shape"/>
		<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
		<feOffset/>
		<feGaussianBlur stdDeviation="11"/>
		<feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1"/>
		<feColorMatrix type="matrix" values="0 0 0 0 0.223529 0 0 0 0 0.85098 0 0 0 0 1 0 0 0 0.25 0"/>
		<feBlend mode="normal" in2="shape" result="effect2_innerShadow_298_336"/>
		</filter>
		<linearGradient id={gradientId} x1="71" y1="159.2" x2="377" y2="17" gradientUnits="userSpaceOnUse">
			<stop stopColor="var(--button-gradient-start)"/>
			<stop offset="0.570524" stopColor="var(--button-gradient-middle)"/>
			<stop offset="1" stopColor="var(--button-gradient-end)"/>
		</linearGradient>
		</defs>
	</svg>
  );
}

export function LargeButtonSvg({ className }: { className?: string }) {
  const uid = useId();
  const gradientId = `primaryLargeGradient-${uid}`;
  const filterId = `primaryLargeFilter-${uid}`;

  return (
	<svg width="534" height="159" viewBox="0 0 534 159" fill="none" className={className}>
		<g filter={`url(#${filterId})`}>
		<path d="M25 25L33 17H501L509 25L517 33V126L509 134L501 142H33L25.3278 134.328L17 126V33L25 25Z" fill={`url(#${gradientId})`}/>
		<path d="M500.586 18L516 33.4141V125.586L500.586 141H33.4141L18 125.586V33.4141L33.4141 18H500.586Z" stroke="#64D0E8" stroke-width="2"/>
		</g>
		<defs>
		<filter id={filterId} x="0" y="0" width="534" height="159" filterUnits="userSpaceOnUse" color-interpolation-filters="sRGB">
		<feFlood flood-opacity="0" result="BackgroundImageFix"/>
		<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
		<feOffset/>
		<feGaussianBlur stdDeviation="8.5"/>
		<feComposite in2="hardAlpha" operator="out"/>
		<feColorMatrix type="matrix" values="0 0 0 0 0.156863 0 0 0 0 0.843137 0 0 0 0 1 0 0 0 0.3 0"/>
		<feBlend mode="normal" in2="BackgroundImageFix" result="effect1_dropShadow_298_334"/>
		<feBlend mode="normal" in="SourceGraphic" in2="effect1_dropShadow_298_334" result="shape"/>
		<feColorMatrix in="SourceAlpha" type="matrix" values="0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0" result="hardAlpha"/>
		<feOffset/>
		<feGaussianBlur stdDeviation="11"/>
		<feComposite in2="hardAlpha" operator="arithmetic" k2="-1" k3="1"/>
		<feColorMatrix type="matrix" values="0 0 0 0 0.223529 0 0 0 0 0.85098 0 0 0 0 1 0 0 0 0.25 0"/>
		<feBlend mode="normal" in2="shape" result="effect2_innerShadow_298_334"/>
		</filter>
		<linearGradient id={gradientId}>
			<stop stopColor="var(--button-gradient-start)"/>
			<stop offset="0.570524" stopColor="var(--button-gradient-middle)"/>
			<stop offset="1" stopColor="var(--button-gradient-end)"/>
		</linearGradient>
		</defs>
	</svg>

  );
}