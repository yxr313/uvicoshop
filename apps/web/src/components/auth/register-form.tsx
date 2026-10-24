"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { ArrowRight, Loader2, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";

const registerSchema = z.object({
  fullName: z.string().min(3, { message: "Ingresa tu nombre completo." }),
  storeSlug: z.string()
    .min(3, { message: "El nombre de tu tienda debe tener al menos 3 caracteres." })
    .regex(/^[a-z0-9-]+$/, { message: "Solo letras minúsculas, números y guiones." }),
  email: z.string().email({ message: "Ingresa un correo válido." }),
  password: z.string()
    .min(8, { message: "Mínimo 8 caracteres." })
    .regex(/[A-Z]/, { message: "Debe tener al menos una mayúscula." })
    .regex(/[0-9]/, { message: "Debe tener al menos un número." }),
  terms: z.boolean().refine((val) => val === true, { message: "Debes aceptar los términos." }),
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export function RegisterForm() {
  const router = useRouter();
  const supabase = createClient();
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { storeSlug: "", terms: false },
  });

  const storeSlug = watch("storeSlug") || "";
  const password = watch("password") || "";
  const termsAccepted = watch("terms");

  // Lógica simple para el medidor de fuerza
  const passwordStrength = {
    length: password.length >= 8,
    upper: /[A-Z]/.test(password),
    number: /[0-9]/.test(password),
  };

  const onSubmit = async (data: RegisterFormValues) => {
    setError(null);
    const { error: authError } = await supabase.auth.signUp({
      email: data.email,
      password: data.password,
      options: {
        data: {
          full_name: data.fullName,
          store_slug: data.storeSlug,
        },
      },
    });

    if (authError) {
      setError(authError.message);
      return;
    }

    router.push("/dashboard");
    router.refresh();
  };

  return (
    <div className="w-full">
      <div className="mb-8">
        <h1 className="text-display-3 font-display text-foreground mb-2">Crea tu tienda.</h1>
        <p className="text-body-sm text-muted-foreground">Sin tarjeta. Configuración en 2 minutos.</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        
        <div className="space-y-1">
          <label className="text-sm font-medium text-foreground">Nombre completo</label>
          <input
            {...register("fullName")}
            type="text"
            placeholder="Ej. Juan Pérez"
            className="w-full h-11 px-4 rounded-xl bg-neutral-50 border border-transparent focus:border-neutral-300 focus:bg-white focus:ring-0 transition-all text-sm outline-none"
          />
          {errors.fullName && <p className="text-[13px] text-destructive mt-1">{errors.fullName.message}</p>}
        </div>

        <div className="space-y-1">
          <label className="text-sm font-medium text-foreground">Nombre de tu tienda</label>
          <div className="flex items-center w-full h-11 rounded-xl bg-neutral-50 border border-transparent focus-within:border-neutral-300 focus-within:bg-white transition-all overflow-hidden px-4">
            <span className="text-muted-foreground text-sm whitespace-nowrap mr-1">uvico.shop/</span>
            <input
              type="text"
              value={storeSlug}
              onChange={(e) => {
                const val = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "");
                setValue("storeSlug", val, { shouldValidate: true });
              }}
              placeholder="mi-tienda"
              className="w-full bg-transparent outline-none text-foreground text-sm placeholder:text-neutral-400"
            />
          </div>
          {errors.storeSlug && <p className="text-[13px] text-destructive mt-1">{errors.storeSlug.message}</p>}
        </div>

        <div className="space-y-1">
          <label className="text-sm font-medium text-foreground">Correo electrónico</label>
          <input
            {...register("email")}
            type="email"
            placeholder="tu@correo.com"
            className="w-full h-11 px-4 rounded-xl bg-neutral-50 border border-transparent focus:border-neutral-300 focus:bg-white focus:ring-0 transition-all text-sm outline-none"
          />
          {errors.email && <p className="text-[13px] text-destructive mt-1">{errors.email.message}</p>}
        </div>

        <div className="space-y-1">
          <label className="text-sm font-medium text-foreground">Contraseña</label>
          <input
            {...register("password")}
            type="password"
            placeholder="Mínimo 8 caracteres"
            className="w-full h-11 px-4 rounded-xl bg-neutral-50 border border-transparent focus:border-neutral-300 focus:bg-white focus:ring-0 transition-all text-sm outline-none"
          />
          
          {/* Medidor de fuerza */}
          {password.length > 0 && (
            <div className="flex flex-col gap-2 mt-2">
              <div className="flex gap-1 h-1">
                <div className={cn("h-full flex-1 rounded-full", passwordStrength.length ? "bg-accent" : "bg-neutral-200")} />
                <div className={cn("h-full flex-1 rounded-full", passwordStrength.upper ? "bg-accent" : "bg-neutral-200")} />
                <div className={cn("h-full flex-1 rounded-full", passwordStrength.number ? "bg-accent" : "bg-neutral-200")} />
              </div>
              <ul className="text-xs space-y-1 text-muted-foreground">
                <li className={cn("flex items-center gap-1", passwordStrength.length && "text-foreground")}>
                  <CheckCircle2 size={12} className={passwordStrength.length ? "text-accent-foreground fill-accent" : "text-neutral-300"} /> Mínimo 8 caracteres
                </li>
                <li className={cn("flex items-center gap-1", passwordStrength.upper && "text-foreground")}>
                  <CheckCircle2 size={12} className={passwordStrength.upper ? "text-accent-foreground fill-accent" : "text-neutral-300"} /> Una letra mayúscula
                </li>
                <li className={cn("flex items-center gap-1", passwordStrength.number && "text-foreground")}>
                  <CheckCircle2 size={12} className={passwordStrength.number ? "text-accent-foreground fill-accent" : "text-neutral-300"} /> Un número
                </li>
              </ul>
            </div>
          )}
          {errors.password && <p className="text-[13px] text-destructive mt-1">{errors.password.message}</p>}
        </div>

        <div className="flex items-start gap-3 mt-6 mb-2">
          <input
            type="checkbox"
            id="terms"
            {...register("terms")}
            className="mt-1 w-4 h-4 rounded border-neutral-300 text-foreground focus:ring-0 cursor-pointer accent-foreground"
          />
          <label htmlFor="terms" className="text-xs text-muted-foreground leading-tight cursor-pointer">
            Acepto los <Link href="/terminos" className="text-foreground underline" target="_blank">Términos y Condiciones</Link> y la <Link href="/privacidad" className="text-foreground underline" target="_blank">Política de Privacidad</Link> de Uvico.
          </label>
        </div>
        {errors.terms && <p className="text-[13px] text-destructive">{errors.terms.message}</p>}

        {error && (
          <div className="p-3 bg-red-50 text-destructive text-sm rounded-lg border border-red-100">
            {error}
          </div>
        )}

        <Button type="submit" className="w-full h-11 mt-2 rounded-full bg-accent text-accent-foreground hover:bg-accent-hover" disabled={isSubmitting || !termsAccepted}>
          {isSubmitting ? <Loader2 className="animate-spin" size={18} /> : (
            <>Crear cuenta <ArrowRight className="ml-2" size={18} /></>
          )}
        </Button>
      </form>

      <p className="mt-8 text-center text-sm text-muted-foreground">
        ¿Ya tienes cuenta? <Link href="/login" className="text-foreground font-medium hover:underline">Inicia sesión</Link>
      </p>
    </div>
  );
}
