import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function NotificationsPage() {
  return (
    <div className="space-y-6 animate-fade-in max-w-2xl">
      <div>
        <h1 className="text-display-3 font-display text-foreground">Notificaciones</h1>
        <p className="text-body text-muted-foreground">Configura tus alertas.</p>
      </div>
      <Card className="rounded-2xl border-border shadow-sm">
        <CardHeader><CardTitle>Preferencias de correo</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <p className="text-sm font-medium">Nuevos pedidos</p>
            <input type="checkbox" defaultChecked className="w-4 h-4 accent-foreground" />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
