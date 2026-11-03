"use client";

import { useEffect, useState } from "react";
import type { UseFormReturn } from "react-hook-form";
import type { ProductFormValues } from "@/hooks/use-product-form";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Link from "@tiptap/extension-link";
import Underline from "@tiptap/extension-underline";
import Placeholder from "@tiptap/extension-placeholder";
import { 
  IconBox, IconTruckDelivery, IconSparkles, IconCheck,
  IconBold, IconItalic, IconUnderline, IconList, IconListNumbers, 
  IconLink, IconArrowBackUp, IconArrowForwardUp, IconH2, IconH3, IconBlockquote
} from "@tabler/icons-react";
import { cn } from "@/lib/utils";

interface Step1Props {
  form: UseFormReturn<ProductFormValues>;
}

// Componente interno para la Barra de Herramientas del Editor
function EditorToolbar({ 
  editor, 
  onGenerateAI, 
  onOpenLinkDialog 
}: { 
  editor: any; 
  onGenerateAI: () => void; 
  onOpenLinkDialog: () => void;
}) {
  if (!editor) return null;

  return (
    <div className="flex flex-wrap items-center gap-1 border-b border-border bg-neutral-50/50 p-2">
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("bold") && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleBold().run()}
        title="Negrita"
      >
        <IconBold className="size-4" />
      </Button>
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("italic") && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleItalic().run()}
        title="Cursiva"
      >
        <IconItalic className="size-4" />
      </Button>
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("underline") && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleUnderline().run()}
        title="Subrayado"
      >
        <IconUnderline className="size-4" />
      </Button>
      
      <div className="w-px h-4 bg-border mx-1" />
      
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("heading", { level: 2 }) && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        title="Título H2"
      >
        <IconH2 className="size-4" />
      </Button>
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("heading", { level: 3 }) && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
        title="Título H3"
      >
        <IconH3 className="size-4" />
      </Button>

      <div className="w-px h-4 bg-border mx-1" />

      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("bulletList") && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleBulletList().run()}
        title="Lista"
      >
        <IconList className="size-4" />
      </Button>
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("orderedList") && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleOrderedList().run()}
        title="Lista numerada"
      >
        <IconListNumbers className="size-4" />
      </Button>
      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("blockquote") && "bg-neutral-200")}
        onClick={() => editor.chain().focus().toggleBlockquote().run()}
        title="Cita"
      >
        <IconBlockquote className="size-4" />
      </Button>

      <div className="w-px h-4 bg-border mx-1" />

      <Button
        type="button" variant="ghost" size="icon" className={cn("size-8", editor.isActive("link") && "bg-neutral-200")}
        onClick={onOpenLinkDialog}
        title="Añadir enlace"
      >
        <IconLink className="size-4" />
      </Button>

      <div className="w-px h-4 bg-border mx-1" />

      <Button
        type="button" variant="ghost" size="icon" className="size-8"
        onClick={() => editor.chain().focus().undo().run()}
        disabled={!editor.can().undo()}
        title="Deshacer"
      >
        <IconArrowBackUp className="size-4" />
      </Button>
      <Button
        type="button" variant="ghost" size="icon" className="size-8"
        onClick={() => editor.chain().focus().redo().run()}
        disabled={!editor.can().redo()}
        title="Rehacer"
      >
        <IconArrowForwardUp className="size-4" />
      </Button>
      
      <div className="ml-auto">
        <Button 
          type="button" 
          variant="ghost" 
          size="sm" 
          onClick={onGenerateAI}
          className="text-xs h-8 text-muted-foreground hover:text-foreground gap-1.5"
        >
          <IconSparkles className="size-3.5 text-[#A7C957]" /> 
          Generar con IA
        </Button>
      </div>
    </div>
  );
}

export function Step1Identity({ form }: Step1Props) {
  const { register, watch, setValue, formState: { errors } } = form;
  const saleMode = watch("sale_mode");
  const name = watch("name");
  const description = watch("description");

  const [isLinkDialogOpen, setIsLinkDialogOpen] = useState(false);
  const [linkUrl, setLinkUrl] = useState("");

  const editor = useEditor({
    immediatelyRender: false,
    extensions: [
      StarterKit,
      Underline,
      Link.configure({ openOnClick: false }),
      Placeholder.configure({
        placeholder: "Describe tu producto. Usa negritas, listas o enlaces para hacerlo atractivo...",
      }),
    ],
    content: description || "",
    onUpdate: ({ editor }) => {
      setValue("description", editor.getHTML(), { shouldValidate: true });
    },
    editorProps: {
      attributes: {
        // AÑADIDOS ESTILOS EXPLÍCITOS PARA LISTAS Y CITAS
        class: "focus:outline-none min-h-[180px] max-w-none px-4 py-3 text-foreground [&_p]:mb-2 [&_ul]:list-disc [&_ul]:pl-5 [&_ul]:mb-2 [&_ol]:list-decimal [&_ol]:pl-5 [&_ol]:mb-2 [&_h2]:text-xl [&_h2]:font-semibold [&_h2]:mb-2 [&_h3]:text-lg [&_h3]:font-semibold [&_h3]:mb-2 [&_blockquote]:border-l-4 [&_blockquote]:border-neutral-300 [&_blockquote]:pl-4 [&_blockquote]:italic [&_a]:text-[#2736D9] [&_a]:underline",
      },
    },
  });

  useEffect(() => {
    if (editor && description && editor.getHTML() !== description) {
      editor.commands.setContent(description);
    }
  }, [editor]);

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setValue("name", val, { shouldValidate: true });
    const generatedSlug = val
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/(^-|-$)+/g, "");
    setValue("slug", generatedSlug, { shouldValidate: true });
  };

  const openLinkDialog = () => {
    if (!editor) return;
    const previousUrl = editor.getAttributes("link").href || "";
    setLinkUrl(previousUrl);
    setIsLinkDialogOpen(true);
  };

  const applyLink = () => {
    if (!editor) return;
    if (linkUrl === "") {
      editor.chain().focus().extendMarkRange("link").unsetLink().run();
    } else {
      const urlWithProtocol = linkUrl.match(/^https?:\/\//) ? linkUrl : `https://${linkUrl}`;
      editor.chain().focus().extendMarkRange("link").setLink({ href: urlWithProtocol }).run();
    }
    setIsLinkDialogOpen(false);
    setLinkUrl("");
  };

  const handleGenerateAI = () => {
    console.log("Llamando a la Groq API...");
  };

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      
      {/* 1. Selector de Modo de Venta */}
      <div className="space-y-4">
        <div>
          <h3 className="text-lg font-medium text-foreground">¿Cómo vas a vender este producto?</h3>
          <p className="text-sm text-muted-foreground">Elige el método de entrega principal. Podrás cambiarlo después.</p>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Tarjeta: Inventario Propio */}
          <Card 
            className={cn(
              "cursor-pointer transition-all relative overflow-hidden",
              saleMode === "own_inventory" 
                ? "border-foreground ring-1 ring-foreground bg-white" 
                : "border-border bg-neutral-50 hover:bg-neutral-100 opacity-70 hover:opacity-100"
            )}
            onClick={() => setValue("sale_mode", "own_inventory", { shouldValidate: true })}
          >
            {/* pointer-events-none asegura que el check no bloquee el clic */}
            <div className={cn(
              "absolute top-0 right-0 bg-foreground text-background p-1.5 rounded-bl-lg transition-opacity pointer-events-none",
              saleMode === "own_inventory" ? "opacity-100" : "opacity-0"
            )}>
              <IconCheck className="size-4" />
            </div>
            <CardContent className="p-5 flex items-start gap-4">
              <div className={cn(
                "p-2.5 rounded-xl shrink-0 transition-colors",
                saleMode === "own_inventory" ? "bg-foreground text-background" : "bg-neutral-200 text-neutral-600"
              )}>
                <IconBox className="size-6" />
              </div>
              <div>
                <h4 className="font-semibold text-foreground text-base">Mi propio inventario</h4>
                <p className="text-sm text-muted-foreground mt-1 leading-relaxed">
                  Yo empaco y envío los productos. Uvico solo te cobra la suscripción.
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Tarjeta: Dropshipping / POD */}
          <Card 
            className={cn(
              "cursor-pointer transition-all relative overflow-hidden",
              saleMode === "dropshipping_verified" 
                ? "border-foreground ring-1 ring-foreground bg-white" 
                : "border-border bg-neutral-50 hover:bg-neutral-100 opacity-70 hover:opacity-100"
            )}
            onClick={() => setValue("sale_mode", "dropshipping_verified", { shouldValidate: true })}
          >
            <div className={cn(
              "absolute top-0 right-0 bg-foreground text-background p-1.5 rounded-bl-lg transition-opacity pointer-events-none",
              saleMode === "dropshipping_verified" ? "opacity-100" : "opacity-0"
            )}>
              <IconCheck className="size-4" />
            </div>
            <CardContent className="p-5 flex items-start gap-4">
              <div className={cn(
                "p-2.5 rounded-xl shrink-0 transition-colors",
                saleMode === "dropshipping_verified" ? "bg-foreground text-background" : "bg-neutral-200 text-neutral-600"
              )}>
                <IconTruckDelivery className="size-6" />
              </div>
              <div>
                <h4 className="font-semibold text-foreground text-base">Dropshipping / POD</h4>
                <p className="text-sm text-muted-foreground mt-1 leading-relaxed">
                  Uvico y sus proveedores se encargan del empaque y envío automático.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* 2. Nombre y Slug */}
      <div className="grid gap-6 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="name">Nombre del producto</Label>
          <Input
            id="name"
            placeholder="Ej. Camiseta Oversize Negra"
            value={name || ""}
            onChange={handleNameChange}
            className="h-11 bg-white" // Forzamos fondo blanco
          />
          {errors.name && <p className="text-xs text-destructive">{errors.name.message}</p>}
        </div>
        
        <div className="space-y-2">
          <Label htmlFor="slug">URL de la tienda</Label>
          <div className="flex items-center h-11 rounded-lg border border-border bg-white focus-within:ring-1 focus-within:ring-neutral-300 focus-within:border-neutral-400 transition-all overflow-hidden">
            <span className="text-sm text-muted-foreground px-3 bg-neutral-50 h-full flex items-center border-r border-border">
              uvico.shop/mi-tienda/
            </span>
            <input
              id="slug"
              className="flex-1 bg-transparent px-3 text-sm font-medium outline-none text-foreground"
              {...register("slug")}
              readOnly
            />
          </div>
          <p className="text-xs text-muted-foreground">Se genera automáticamente a partir del nombre.</p>
        </div>
      </div>

      {/* 3. Editor de Texto Enriquecido (TipTap Completo) */}
      <div className="space-y-2">
        <Label>Descripción del producto</Label>
        <div className="border border-border rounded-lg overflow-hidden bg-white focus-within:border-neutral-400 transition-all">
          <EditorToolbar 
            editor={editor} 
            onGenerateAI={handleGenerateAI} 
            onOpenLinkDialog={openLinkDialog}
          />
          <EditorContent editor={editor} />
        </div>
        {errors.description && <p className="text-xs text-destructive">{errors.description.message}</p>}
      </div>

      {/* Dialog de Link */}
      <Dialog open={isLinkDialogOpen} onOpenChange={setIsLinkDialogOpen}>
        {/* Forzamos bg-white y border-border en el DialogContent */}
        <DialogContent className="sm:max-w-sm bg-white border-border shadow-lg">
          <DialogHeader>
            <DialogTitle>Añadir Enlace</DialogTitle>
          </DialogHeader>
          <div className="py-4">
            <Label>URL del enlace</Label>
            {/* Forzamos bg-white y border-border en el Input */}
            <Input 
              placeholder="ej. https://mi-sitio.com" 
              value={linkUrl} 
              onChange={(e) => setLinkUrl(e.target.value)} 
              className="mt-2 bg-white border-border text-foreground"
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setIsLinkDialogOpen(false)}>Cancelar</Button>
            <Button onClick={applyLink} className="bg-foreground text-background hover:bg-foreground/90">Aplicar</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}