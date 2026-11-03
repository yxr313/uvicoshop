import { useEffect } from "react";
import { useForm, type UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";

export const productSchema = z.object({
  // Paso 1
  name: z.string().min(3, "El nombre debe tener al menos 3 caracteres"),
  slug: z.string().min(3, "El slug es requerido"),
  description: z.string().optional(),
  sale_mode: z.enum(["own_inventory", "dropshipping_verified"]),
  
  // Paso 2
  price_cents: z.coerce.number().min(1, "El precio es requerido"),
  cost_cents: z.coerce.number().min(0).optional(),
  shipping_cost_cents: z.coerce.number().min(0).optional(), // NUEVO
  ad_cost_cents: z.coerce.number().min(0).optional(),       // NUEVO
  initial_stock: z.coerce.number().min(0).optional(),
  provider_id: z.string().optional(),

  // Paso 3
  images: z.array(z.string()).optional(),
  has_variations: z.boolean().optional(),
  variations: z.array(z.object({
    name: z.string(),
    options: z.array(z.string())
  })).optional(),

  // Paso 4
  status: z.enum(["draft", "published"]).optional(),
  seo_title: z.string().optional(),
  seo_description: z.string().optional(),
  tags: z.string().optional(),
});

export type ProductFormValues = z.infer<typeof productSchema>;

const DRAFT_KEY = "uvico-product-draft";

export function useProductForm(): { form: UseFormReturn<ProductFormValues>; clearDraft: () => void } {
  const getDraft = (): Partial<ProductFormValues> => {
    if (typeof window === "undefined") return {};
    try {
      const saved = localStorage.getItem(DRAFT_KEY);
      return saved ? (JSON.parse(saved) as Partial<ProductFormValues>) : {};
    } catch {
      return {};
    }
  };

  const form = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      name: "",
      slug: "",
      description: "",
      sale_mode: "own_inventory",
      price_cents: 0,
      cost_cents: 0,
      shipping_cost_cents: 0, // NUEVO
      ad_cost_cents: 0,       // NUEVO
      initial_stock: 0,
      provider_id: "",
      images: [],
      has_variations: false,
      variations: [],
      status: "draft",
      seo_title: "",
      seo_description: "",
      tags: "",
      ...getDraft(),
    },
  });

  useEffect(() => {
    const subscription = form.watch((value) => {
      localStorage.setItem(DRAFT_KEY, JSON.stringify(value));
    });
    return () => subscription.unsubscribe();
  }, [form]);

  const clearDraft = () => localStorage.removeItem(DRAFT_KEY);

  return { form, clearDraft };
}