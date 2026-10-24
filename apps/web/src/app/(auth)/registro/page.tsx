"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { ArrowRight, Loader2, Mail } from "lucide-react";

const registerSchema = z.object({
  fullName: z.string().min(3, { message: "Ingresa tu nombre completo." }),
  storeSlug: z.string()
    .min(3, { message: "El nombre de tu tienda debe tener al menos 3 caracteres." })
    .regex(/^[a-z0-9-]+$/, { message: "Solo letras minúsculas, números y guiones." }),
  email: z.string().email({ message: "Ingresa un correo válido." }),
  password: z.string().min(8, { message: "La contraseña debe tener al menos 8 caracteres." }),
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export default function RegisterPage() {
  const router = useRouter();
  const supabase = createClient();
  const [error, setError] = useState<string | null>(null);
  
  // Estados para el flujo de Vercel: 'form' o 'verify'
  const [step, setStep] = useState<"form" | "verify">("form");
  const [registeredEmail, setRegisteredEmail] = useState("");

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
  });

  const storeSlug = watch("storeSlug");

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

    // Si todo sale bien, pasamos al estado de verificación
    setRegisteredEmail(data.email);
    setStep("verify");
  };

  // --- ESTADO 1: FORMULARIO ---
  if (step === "form") {
    return (
      <div className="w-full animate-fade-in">
        <div className="mb-8">
          <h1 className="text-display-3 font-display text-foreground mb-2">Crea tu cuenta.</h1>
          <p className="text-body-sm text-muted-foreground">
            Únete a los mejores emprendedores de Venezuela.
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Nombre completo</label>
            <input
              {...register("fullName")}
              type="text"
              placeholder="Ej. Juan Pérez"
              className="w-full h-11 px-3 rounded-lg bg-white border border-border focus:border-neutral-400 focus:ring-1 focus:ring-neutral-400 transition-all text-sm outline-none"
            />
            {errors.fullName && <p className="text-xs text-destructive mt-1">{errors.fullName.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Nombre de tu tienda</label>
            <div className="flex items-center w-full h-11 rounded-lg bg-white border border-border focus-within:border-neutral-400 focus-within:ring-1 focus-within:ring-neutral-400 transition-all overflow-hidden px-3">
              <span className="text-muted-foreground text-sm whitespace-nowrap mr-1">uvico.shop/</span>
              <input
                type="text"
                value={storeSlug || ""}
                onChange={(e) => {
                  const val = e.target.value.toLowerCase().replace(/[^a-z0-9-]/g, "");
                  setValue("storeSlug", val, { shouldValidate: true });
                }}
                placeholder="mi-tienda"
                className="w-full bg-transparent outline-none text-foreground text-sm placeholder:text-neutral-400"
              />
            </div>
            {errors.storeSlug && <p className="text-xs text-destructive mt-1">{errors.storeSlug.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Correo electrónico</label>
            <input
              {...register("email")}
              type="email"
              placeholder="tu@correo.com"
              className="w-full h-11 px-3 rounded-lg bg-white border border-border focus:border-neutral-400 focus:ring-1 focus:ring-neutral-400 transition-all text-sm outline-none"
            />
            {errors.email && <p className="text-xs text-destructive mt-1">{errors.email.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">Contraseña</label>
            <input
              {...register("password")}
              type="password"
              placeholder="Mínimo 8 caracteres"
              className="w-full h-11 px-3 rounded-lg bg-white border border-border focus:border-neutral-400 focus:ring-1 focus:ring-neutral-400 transition-all text-sm outline-none"
            />
            {errors.password && <p className="text-xs text-destructive mt-1">{errors.password.message}</p>}
          </div>

          {error && (
            <div className="p-3 bg-red-50 text-destructive text-sm rounded-lg border border-red-100">
              {error}
            </div>
          )}

          <Button type="submit" className="w-full h-11 mt-6 bg-accent text-accent-foreground hover:bg-accent-hover rounded-lg" disabled={isSubmitting}>
            {isSubmitting ? <Loader2 className="animate-spin" size={18} /> : (
              <>Crear cuenta <ArrowRight className="ml-2" size={18} /></>
            )}
          </Button>
        </form>

        <p className="mt-8 text-sm text-muted-foreground">
          ¿Ya tienes cuenta? <Link href="/login" className="text-foreground font-medium hover:underline">Inicia sesión</Link>
        </p>
        
        <p className="mt-8 text-xs text-muted-foreground leading-relaxed">
          Al registrarte, aceptas nuestros <Link href="/terminos" className="underline">Términos de Servicio</Link> y <Link href="/privacidad" className="underline">Política de Privacidad</Link>.
        </p>
      </div>
    );
  }

  // --- ESTADO 2: VERIFICACIÓN (Estilo Vercel) ---
  return (
    <div className="w-full text-center animate-fade-in">
      <div className="w-12 h-12 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-6">
        <Mail className="w-6 h-6 text-foreground" />
      </div>
      
      <h1 className="text-display-3 font-display text-foreground mb-4">Revisa tu correo.</h1>
      
      <p className="text-body text-muted-foreground mb-8">
        Hemos enviado un enlace de verificación a <br/>
        <strong className="text-foreground">{registeredEmail}</strong>
      </p>

      <div className="bg-neutral-50 border border-border p-4 rounded-xl mb-8">
        <p className="text-sm text-muted-foreground">
          Haz clic en el enlace del correo para confirmar tu cuenta y acceder a tu panel de control. Si no lo ves, revisa tu carpeta de spam.
        </p>
      </div>

      <Link href="/login">
        <Button variant="outline" className="w-full h-11 rounded-lg">
          Ir a Iniciar Sesión
        </Button>
      </Link>
    </div>
  );
}
