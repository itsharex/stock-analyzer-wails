import React from 'react'
import { IndustryInfo } from '../types'
import { Building2, Tag } from 'lucide-react'

interface IndustryPanelProps {
  industryInfo: IndustryInfo | undefined
}

const IndustryPanel: React.FC<IndustryPanelProps> = ({ industryInfo }) => {
  if (!industryInfo) {
    return (
      <div className="p-4 bg-gray-800 rounded-lg shadow-lg">
        <h3 className="text-lg font-semibold text-gray-300 mb-2 flex items-center">
          <Building2 className="w-5 h-5 mr-2 text-amber-400" /> 行业宏观
        </h3>
        <p className="text-gray-500">暂无行业数据</p>
      </div>
    )
  }

  return (
    <div className="p-4 bg-gray-800 rounded-lg shadow-lg">
      <h3 className="text-lg font-semibold text-gray-300 mb-2 flex items-center">
        <Building2 className="w-5 h-5 mr-2 text-amber-400" /> 行业宏观
      </h3>
      <div className="space-y-3 text-sm">
        <div className="flex items-center">
          <span className="text-gray-400 w-20">所属行业:</span>
          <span className="font-medium text-white">{industryInfo.industry_name}</span>
        </div>
        <div className="flex items-center">
          <span className="text-gray-400 w-20">行业 PE:</span>
          <span className="font-medium text-white">{industryInfo.industry_pe.toFixed(2)}x</span>
        </div>
        <div className="flex items-start">
          <span className="text-gray-400 w-20 mt-1">概念板块:</span>
          <div className="flex flex-wrap gap-1">
            {industryInfo.concept_names.map((tag, index) => (
              <span key={index} className="px-2 py-0.5 text-xs bg-indigo-700 text-indigo-200 rounded-full flex items-center">
                <Tag className="w-3 h-3 mr-1" /> {tag}
              </span>
            ))}
          </div>
        </div>
        
      </div>
    </div>
  )
}

export default IndustryPanel
