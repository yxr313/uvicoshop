"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { Loader2, CheckCircle2, XCircle } from "lucide-react";

interface SlugBuilderProps {
  variant?: "default" | "dark";
  className?: string;
  initialSlug?: string;
}

export function SlugBuilder({ variant = "dark", className, initialSlug = "" }: SlugBuilderProps) {
  const router = useRouter();
  const [slug, setSlug] = useState(initialSlug);
  const [isChecking, setIsChecking] = useState(false);
  const [isAvailable, setIsAvailable] = useState<boolean | null>(null);

  useEffect(() => {
    if (slug.length < 3) {
      setIsAvailable(null);
      return;
    }

    const timer = setTimeout(async () => {
      setIsChecking(true);
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/v1";
        const res = await fetch(`${apiUrl}/tenants/check-slug?slug=${slug}`);
        if (res.ok) {
          const data = await res.json();
          setIsAvailable(data.available);
        } else {
          setIsAvailable(null);
        }
      } catch (error) {
        console.error("Error checking slug:", error);
        setIsAvailable(null);
      } finally {
        setIsChecking(false);
      }
    }, 500);

    return () => clearTimeout(timer);
  }, [slug]);

  const handleSlugChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "");
    setSlug(val);
  };

  const handleCreate = () => {
    if (isAvailable && !isChecking) {
      router.push(`/registro?slug=${slug}`);
    }
  };

  const getButtonText = () => {
    if (slug.length < 3) return "Crear tienda";
    if (isChecking) return "Verificando...";
    if (isAvailable === false) return "No disponible";
    return "Crear tienda";
  };

  const isButtonDisabled = isAvailable !== true || isChecking;

  // Mapeo seguro: "dark" → "secondary", "default" → "default"
  const buttonVariant = variant === "dark" ? "secondary" : "default";

  return (
    <div className={cn("w-full", className)}>
      <div
        className={cn(
          "flex flex-col sm:flex-row items-center w-full p-1.5 bg-white border rounded-full shadow-sm transition-all duration-300",
          "focus-within:shadow-md",
          isAvailable === false ? "border-red-400 focus-within:border-red-400" :
          isAvailable === true ? "border-green-400 focus-within:border-green-400" :
          "border-border focus-within:border-neutral-300"
        )}
      >
        <div className="flex items-center flex-1 w-full px-4 h-10">
          <span className="text-muted-foreground font-medium whitespace-nowrap mr-1 text-sm">
            uvico.shop/
          </span>
          <input
            type="text"
            value={slug}
            onChange={handleSlugChange}
            placeholder="mi-tienda"
            className="w-full bg-transparent outline-none placeholder:text-neutral-300 text-foreground font-medium text-sm"
          />
          <div className="flex items-center justify-center w-6 h-6 ml-2 shrink-0">
            {isChecking && <Loader2 className="w-4 h-4 text-muted-foreground animate-spin" />}
            {!isChecking && isAvailable === true && <CheckCircle2 className="w-4 h-4 text-green-500" />}
            {!isChecking && isAvailable === false && <XCircle className="w-4 h-4 text-red-500" />}
          </div>
        </div>

        <Button
          variant={buttonVariant}
          onClick={handleCreate}
          disabled={isButtonDisabled}
          className={cn(
            "w-full sm:w-auto mt-1.5 sm:mt-0 h-10 px-6 rounded-full text-sm transition-all duration-300 min-w-[130px]",
            isButtonDisabled && "opacity-50 cursor-not-allowed",
            isAvailable === false && "bg-red-50 text-red-600 border-red-200 hover:bg-red-50"
          )}
        >
          {getButtonText()}
        </Button>
      </div>
    </div>
  );
}