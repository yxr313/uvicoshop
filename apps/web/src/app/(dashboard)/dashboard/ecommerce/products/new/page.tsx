"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useForm, type Path } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { PageHeader } from "@/components/page-header";
import { Button } from "@/components/ui/button";
import { IconArrowLeft, IconArrowRight, IconDeviceFloppy, IconCheck, IconLoader2 } from "@tabler/icons-react";
import { cn } from "@/lib/utils";
import { fetchFromApi } from "@/lib/api-client";
import { useProductForm, type ProductFormValues } from "@/hooks/use-product-form";
import { Step1Identity } from "@/components/app/product-form/step-1-identity";
import { Step2Pricing } from "@/components/app/product-form/step-2-pricing";
import { Step3Gallery } from "@/components/app/product-form/step-3-gallery";
import { Step4Publish } from "@/components/app/product-form/step-4-publish";

const STEPS = [
  { id: 1, title: "Identidad" },
  { id: 2, title: "Precios" },
  { id: 3, title: "Galería" },
  { id: 4, title: "Publicación" },
];

interface CreateProductPayload {
  sale_mode: string;
  provider_id?: string;
  sku: string;
  slug: string;
  name: string;
  description: string;
  price_cents: number;
  initial_stock: number;
}

interface ProductResponse {
  id: string;
}

export default function NewProductPage() {
  const router = useRouter();
  const [currentStep, setCurrentStep] = useState(1);
  const { form, clearDraft } = useProductForm();

  const createProductMutation = useMutation({
    mutationFn: async (data: ProductFormValues) => {
      const generatedSku = `${data.name.toLowerCase().replace(/[^a-z0-9]/g, "-")}-${Math.floor(Math.random() * 10000)}`;

      const payload: CreateProductPayload = {
        sale_mode: data.sale_mode,
        provider_id: data.provider_id || undefined,
        sku: generatedSku,
        slug: data.slug,
        name: data.name,
        description: data.description || "",
        price_cents: data.price_cents,
        initial_stock: data.sale_mode === "own_inventory" ? (data.initial_stock || 0) : 0,
      };

      const product = await fetchFromApi<ProductResponse>("/products", {
        method: "POST",
        body: JSON.stringify(payload),
      });

      if (data.status === "published") {
        await fetchFromApi(`/products/${product.id}/publish`, {
          method: "POST",
        });
      }

      return product;
    },
    onSuccess: () => {
      clearDraft();
      toast.success("Producto creado exitosamente", {
        description: "Tu producto ya está en el catálogo.",
      });
      router.push("/dashboard/ecommerce/products");
    },
    onError: (error: any) => {
      toast.error("Error al crear el producto", {
        description: error.message || "Verifica los datos e intenta de nuevo.",
      });
    },
  });

  const handleNext = async () => {
    let fieldsToValidate: Path<ProductFormValues>[] = [];
    if (currentStep === 1) fieldsToValidate = ["name", "slug", "sale_mode"];
    if (currentStep === 2) fieldsToValidate = ["price_cents"];

    const isValid = await form.trigger(fieldsToValidate);
    if (isValid) setCurrentStep((s) => Math.min(s + 1, 4));
  };

  const handlePublish = async () => {
    const isValid = await form.trigger();
    if (!isValid) {
      toast.error("Faltan datos", {
        description: "Revisa los pasos anteriores antes de publicar.",
      });
      return;
    }

    const data = form.getValues();
    createProductMutation.mutate(data);
  };

  const handleSaveDraft = async () => {
    toast.success("Borrador guardado", {
      description: "Puedes cerrar esta pestaña y continuar después.",
    });
  };

  return (
    <div className="space-y-8 max-w-5xl mx-auto">
      <PageHeader
        title="Crear Producto"
        description="Sigue los pasos para añadir un producto a tu catálogo."
      >
        <Button variant="outline" onClick={() => router.push("/dashboard/ecommerce/products")}>
          <IconArrowLeft className="size-4 mr-2" /> Volver
        </Button>
      </PageHeader>

      {/* Stepper Premium Circular */}
      <div className="flex items-center justify-between w-full px-4 mb-8">
        {STEPS.map((step, index) => (
          <div key={step.id} className="flex items-center flex-1 last:flex-none">
            <div className="flex flex-col items-center">
              <div
                className={cn(
                  "size-10 rounded-full flex items-center justify-center text-sm font-medium border-2 transition-all duration-300",
                  currentStep > step.id
                    ? "bg-foreground border-foreground text-background"
                    : currentStep === step.id
                    ? "border-foreground text-foreground ring-4 ring-neutral-100"
                    : "border-border text-muted-foreground"
                )}
              >
                {currentStep > step.id ? <IconCheck className="size-5" /> : step.id}
              </div>
              <span
                className={cn(
                  "text-xs mt-3 font-medium",
                  currentStep >= step.id ? "text-foreground" : "text-muted-foreground"
                )}
              >
                {step.title}
              </span>
            </div>
            {index < STEPS.length - 1 && (
              <div
                className={cn(
                  "h-[2px] flex-1 mx-4 transition-colors duration-300",
                  currentStep > step.id ? "bg-foreground" : "bg-border"
                )}
              />
            )}
          </div>
        ))}
      </div>

      {/* Formulario Flotante */}
      <div className="bg-white p-6 md:p-8 rounded-3xl border border-border shadow-sm min-h-[400px]">
        <form className="space-y-8" onSubmit={(e) => e.preventDefault()}>
          {currentStep === 1 && <Step1Identity form={form} />}
          {currentStep === 2 && <Step2Pricing form={form} />}
          {currentStep === 3 && <Step3Gallery form={form} />}
          {currentStep === 4 && <Step4Publish form={form} />}
        </form>
      </div>

      {/* Barra de Acciones Inferior (Estática, alineada al contenedor) */}
      <div className="flex justify-between items-center pt-2 pb-8">
        <Button
          variant="outline"
          disabled={currentStep === 1 || createProductMutation.isPending}
          onClick={() => setCurrentStep((s) => Math.max(s - 1, 1))}
        >
          Atrás
        </Button>

        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            onClick={handleSaveDraft}
            disabled={createProductMutation.isPending}
            className="text-muted-foreground hover:text-foreground"
          >
            <IconDeviceFloppy className="size-4 mr-2" /> Guardar Borrador
          </Button>

          {currentStep < 4 ? (
            <Button onClick={handleNext} className="bg-foreground text-background min-w-[140px]">
              Siguiente <IconArrowRight className="size-4 ml-2" />
            </Button>
          ) : (
            <Button
              onClick={form.handleSubmit(handlePublish)}
              disabled={createProductMutation.isPending}
              className="bg-accent text-accent-foreground hover:bg-accent-hover min-w-[160px]"
            >
              {createProductMutation.isPending ? (
                <>
                  <IconLoader2 className="size-4 mr-2 animate-spin" /> Publicando...
                </>
              ) : (
                <>Publicar Producto</>
              )}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}