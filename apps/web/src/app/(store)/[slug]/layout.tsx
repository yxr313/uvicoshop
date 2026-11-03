"use client";

import { useEffect } from "react";
import { usePathname } from "next/navigation";

export function useTracker(tenantId: string) {
  const pathname = usePathname();

  useEffect(() => {
    // 1. Obtener o generar un ID de sesión
    let sessionId = sessionStorage.getItem("uvico_session");
    if (!sessionId) {
      sessionId = crypto.randomUUID();
      sessionStorage.setItem("uvico_session", sessionId);
    }

    // 2. Detectar dispositivo, navegador y OS (Simple)
    const ua = navigator.userAgent;
    const isMobile = /Mobile|Android|iP(hone|od|ad)/.test(ua);
    const device = isMobile ? "Mobile" : "Desktop";
    const browser = ua.includes("Chrome") ? "Chrome" : ua.includes("Safari") ? "Safari" : "Other";
    const os = ua.includes("Win") ? "Windows" : ua.includes("Mac") ? "MacOS" : ua.includes("Android") ? "Android" : "iOS";

    // 3. Obtener UTM de la URL
    const searchParams = new URLSearchParams(window.location.search);
    const utmSource = searchParams.get("utm_source") || "";

    // 4. Enviar al backend en Go
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/analytics/track`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        tenant_id: tenantId,
        session_id: sessionId,
        path: pathname,
        device,
        browser,
        os,
        referrer: document.referrer,
        utm_source: utmSource,
      }),
    }).catch(() => {});
    
  }, [tenantId, pathname]);
}