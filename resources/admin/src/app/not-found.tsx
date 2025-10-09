"use client";

import { useRouter } from "next/navigation";


export default function E404() {
  
  const { push } = useRouter();

  push('/content-manager');
}