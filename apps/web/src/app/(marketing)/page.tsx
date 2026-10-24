import { Hero } from "@/components/marketing/hero/hero";
import { BentoGrid } from "@/components/marketing/showcase/bento-grid";
import { DropshippingCalculator } from "@/components/marketing/calculator/dropshipping-calculator";
import { FAQ } from "@/components/marketing/faq/faq-accordion";
import { FinalCTA } from "@/components/marketing/cta/final-cta";

export default function LandingPage() {
  return (
    <>
      <Hero />
      <BentoGrid />
      <DropshippingCalculator />
      <FAQ />
      <FinalCTA />
    </>
  );
}
