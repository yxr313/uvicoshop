"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { ArrowRight, Loader2 } from "lucide-react";

const loginSchema = z.object({
  email: z.string().email({ message: "Ingresa un correo válido." }),
  password: z.string().min(1, { message: "La contraseña es obligatoria." }),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export function LoginForm() {
  const router = useRouter();
  const supabase = createClient();
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormValues) => {
    setError(null);
    const { error } = await supabase.auth.signInWithPassword({
      email: data.email,
      password: data.password,
    });

    if (error) {
      if (error.message.includes("Invalid login credentials")) {
        setError("El correo o la contraseña son incorrectos.");
      } else if (error.message.includes("Email not confirmed")) {
        setError("Por favor, confirma tu correo electrónico antes de iniciar sesión.");
      } else {
        setError(error.message);
      }
      return;
    }

    router.push("/dashboard");
    router.refresh();
  };

  return (
    <div className="w-full">
      <div className="mb-8">
        <h1 className="text-display-3 font-display text-foreground mb-2">Inicia sesión.</h1>
        <p className="text-body-sm text-muted-foreground">Bienvenido de vuelta a tu panel.</p>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        
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
          <div className="flex justify-between items-center">
            <label className="text-sm font-medium text-foreground">Contraseña</label>
            <Link href="/recuperar" className="text-xs text-muted-foreground hover:text-foreground transition-colors">
              ¿Olvidaste tu contraseña?
            </Link>
          </div>
          <input
            {...register("password")}
            type="password"
            placeholder="••••••••"
            className="w-full h-11 px-4 rounded-xl bg-neutral-50 border border-transparent focus:border-neutral-300 focus:bg-white focus:ring-0 transition-all text-sm outline-none"
          />
          {errors.password && <p className="text-[13px] text-destructive mt-1">{errors.password.message}</p>}
        </div>

        {error && (
          <div className="p-3 bg-red-50 text-destructive text-sm rounded-lg border border-red-100">
            {error}
          </div>
        )}

        <Button type="submit" className="w-full h-11 mt-6 rounded-full" disabled={isSubmitting}>
          {isSubmitting ? <Loader2 className="animate-spin" size={18} /> : (
            <>Iniciar sesión <ArrowRight className="ml-2" size={18} /></>
          )}
        </Button>
      </form>

      <p className="mt-8 text-center text-sm text-muted-foreground">
        ¿No tienes cuenta? <Link href="/registro" className="text-foreground font-medium hover:underline">Crear una tienda</Link>
      </p>
    </div>
  );
}
