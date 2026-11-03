"use client";

import type { UseFormReturn } from "react-hook-form";
import type { ProductFormValues } from "@/hooks/use-product-form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

interface Step4Props {
  form: UseFormReturn<ProductFormValues>;
}

export function Step4Publish({ form }: Step4Props) {
  const { register, watch, setValue } = form;
  const name = watch("name");
  const status = watch("status") || "draft";

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-500">
      
      <div className="space-y-2">
        <Label>Estado de publicación</Label>
        <Select value={status} onValueChange={(val) => setValue("status", val as "draft" | "published", { shouldValidate: true })}>
          <SelectTrigger className="w-full sm:w-[300px]">
            <SelectValue placeholder="Selecciona el estado" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="draft">Guardar como Borrador (Oculto)</SelectItem>
            <SelectItem value="published">Publicar en mi tienda</SelectItem>
          </SelectContent>
        </Select>
        <p className="text-xs text-muted-foreground">Los borradores no son visibles para tus clientes.</p>
      </div>

      <div className="space-y-4 border-t border-border pt-6">
        <h3 className="font-medium text-foreground">Optimización para Google (SEO)</h3>
        
        <div className="space-y-2">
          <Label>Título SEO</Label>
          <Input 
            {...register("seo_title")} 
            placeholder={name ? `Comprar ${name} en Uvico` : "Título para los motores de búsqueda"} 
          />
          <p className="text-xs text-muted-foreground">Recomendado: Máximo 60 caracteres.</p>
        </div>

        <div className="space-y-2">
          <Label>Meta Descripción</Label>
          <Textarea 
            {...register("seo_description")} 
            placeholder="Una breve descripción de tu producto para Google..."
            className="resize-none"
            rows={3}
          />
        </div>
      </div>

      <div className="space-y-2 border-t border-border pt-6">
        <Label>Etiquetas (Tags)</Label>
        <Input 
          {...register("tags")} 
          placeholder="Ej. ropa, algodón, oferta (separadas por coma)" 
        />
        <p className="text-xs text-muted-foreground">Ayudan a los clientes a encontrar productos similares en tu tienda.</p>
      </div>

    </div>
  );
}