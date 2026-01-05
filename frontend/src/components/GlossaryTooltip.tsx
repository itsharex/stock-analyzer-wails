import React, { useState } from 'react';
import { STOCK_GLOSSARY } from '../utils/glossary';
import { HelpCircle, Info, BookOpen } from 'lucide-react';

interface GlossaryTooltipProps {
  term: string;
  children: React.ReactNode;
}

export function GlossaryTooltip({ term, children }: GlossaryTooltipProps) {
  const [isVisible, setIsVisible] = useState(false);
  const item = STOCK_GLOSSARY[term];

  if (!item) return <>{children}</>;

  return (
    <span 
      className="relative inline-flex items-center group cursor-help border-b border-dotted border-blue-400/50 hover:border-blue-400"
      onMouseEnter={() => setIsVisible(true)}
      onMouseLeave={() => setIsVisible(false)}
    >
      {children}
      <HelpCircle className="w-3 h-3 ml-0.5 text-blue-400/70 group-hover:text-blue-400" />
      
      {isVisible && (
        <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 w-64 p-3 bg-slate-800 border border-slate-700 rounded-xl shadow-2xl z-50 animate-in fade-in slide-in-from-bottom-1">
          <div className="flex items-start space-x-2 mb-2">
            <div className="p-1 bg-blue-500/20 rounded">
              <Info className="w-4 h-4 text-blue-400" />
            </div>
            <h4 className="font-bold text-slate-100 text-sm">{item.term}</h4>
          </div>
          <div className="space-y-2">
            <div>
              <p className="text-[10px] text-slate-500 uppercase tracking-wider">专业定义</p>
              <p className="text-xs text-slate-300 leading-relaxed">{item.explanation}</p>
            </div>
            <div className="pt-2 border-t border-slate-700/50">
              <p className="text-[10px] text-blue-400 uppercase tracking-wider font-bold">白话解释</p>
              <p className="text-xs text-blue-100 leading-relaxed italic">“{item.simpleExplanation}”</p>
            </div>
          </div>
          <div className="absolute top-full left-1/2 -translate-x-1/2 border-8 border-transparent border-t-slate-800" />
        </div>
      )}
    </span>
  );
}

export function GlossaryPanel({ text }: { text: string }) {
  const terms = Object.keys(STOCK_GLOSSARY).filter(term => text.includes(term));
  
  if (terms.length === 0) return null;

  return (
    <div className="mt-8 pt-8 border-t border-slate-800">
      <h4 className="text-xs font-bold text-slate-500 uppercase tracking-widest mb-5 flex items-center">
        <BookOpen className="w-4 h-4 mr-2 text-blue-500" />
        报告术语百科
      </h4>
      <div className="grid gap-4">
        {terms.map(term => (
          <div key={term} className="bg-slate-800/40 rounded-xl p-4 border border-slate-700/50 hover:border-blue-500/30 transition-colors">
            <div className="flex items-center space-x-2 mb-2">
              <div className="w-1.5 h-1.5 rounded-full bg-blue-500" />
              <span className="text-sm font-bold text-slate-200">{term}</span>
            </div>
            <p className="text-xs text-slate-400 leading-relaxed pl-3.5 border-l border-slate-700">
              {STOCK_GLOSSARY[term].simpleExplanation}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
