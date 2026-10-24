export function TrustBar() {
  // Simulamos la tipografía de cada logo hasta obtener los SVG
  const gateways = [
    { name: "Wayu Pay", style: "font-sans font-bold tracking-tight" },
    { name: "VeriPay", style: "font-display font-black italic" },
    { name: "VEXPay", style: "font-mono font-medium uppercase tracking-widest" },
    { name: "Payku", style: "font-sans font-light tracking-wide" },
    { name: "Pabilo", style: "font-display font-semibold" },
    { name: "Binance", style: "font-sans font-bold text-[#F3BA2F]" },
    { name: "Zelle", style: "font-serif font-bold text-[#6D1ED4]" },
  ];

  const LogoList = ({ ariaHidden = false }: { ariaHidden?: boolean }) => (
    <ul
      className="flex items-center justify-center md:justify-start [&_li]:mx-8 animate-infinite-scroll min-w-full shrink-0"
      aria-hidden={ariaHidden}
    >
      {gateways.map((gateway, index) => (
        <li key={`${gateway.name}-${index}`} className="flex items-center justify-center">
          <span className={`text-2xl md:text-3xl text-foreground/40 grayscale hover:grayscale-0 hover:text-foreground transition-all duration-300 ${gateway.style}`}>
            {gateway.name}
          </span>
        </li>
      ))}
    </ul>
  );

  return (
    <div className="flex flex-col items-center gap-8 mt-24 w-full">
      <div className="w-full inline-flex flex-nowrap overflow-hidden [mask-image:_linear-gradient(to_right,transparent_0,_black_128px,_black_calc(100%-128px),transparent_100%)]">
        <LogoList />
        <LogoList ariaHidden={true} />
      </div>
    </div>
  );
}
