"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { createClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Loader2, Mail } from "lucide-react";

const recoverSchema = z.object({
  email: z.string().email({ message: "Ingresa un correo válido." }),
});

type RecoverFormValues = z.infer<typeof recoverSchema>;

export function RecoverForm() {
  const supabase = createClient();
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RecoverFormValues>({
    resolver: zodResolver(recoverSchema),
  });

  const onSubmit = async (data: RecoverFormValues) => {
    setError(null);
    const { error } = await supabase.auth.resetPasswordForEmail(data.email, {
      redirectTo: `${window.location.origin}/actualizar-password`,
    });

    if (error) {
      setError(error.message);
      return;
    }

    setSuccess(true);
  };

  if (success) {
    return (
      <div className="w-full text-center animate-fade-in">
        <div className="w-12 h-12 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-6">
          <Mail className="w-6 h-6 text-foreground" />
        </div>
        <h1 className="text-display-3 font-display text-foreground mb-4">Revisa tu correo.</h1>
        <p className="text-body text-muted-foreground mb-8">
          Hemos enviado un enlace para restablecer tu contraseña. Haz clic en él para continuar.
        </p>
        <Link href="/login">
          <Button variant="outline" className="w-full h-11 rounded-full">
            Volver a Iniciar Sesión
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="w-full">
      <div className="mb-8">
        <Link href="/login" className="inline-flex items-center text-sm font-medium text-muted-foreground hover:text-foreground mb-6 transition-colors">
          <ArrowLeft className="w-4 h-4 mr-2" /> Volver
        </Link>
        <h1 className="text-display-3 font-display text-foreground mb-2">Recupera tu acceso.</h1>
        <p className="text-body-sm text-muted-foreground">Te enviaremos un enlace a tu correo.</p>
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

        {error && (
          <div className="p-3 bg-red-50 text-destructive text-sm rounded-lg border border-red-100">
            {error}
          </div>
        )}

        <Button type="submit" className="w-full h-11 mt-6 rounded-full" disabled={isSubmitting}>
          {isSubmitting ? <Loader2 className="animate-spin" size={18} /> : (
            <>Enviar enlace de recuperación</>
          )}
        </Button>
      </form>
    </div>
  );
}
