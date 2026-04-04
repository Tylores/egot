#!/usr/bin/env python3
"""Generate detailed mapping document showing WADL specs vs handler code."""

import xml.etree.ElementTree as ET
import re
from pathlib import Path
from typing import Dict, List


class WADLParser:
    NAMESPACES = {
        'wadl': 'http://wadl.dev.java.net/2009/02',
        'sep': 'urn:ieee:std:2030.5:ns',
        'xsd': 'http://www.w3.org/2001/XMLSchema',
        'wx': 'urn:ieee:std:2030.5:wadlExt'
    }
    
    def __init__(self, wadl_file: str):
        self.wadl_file = wadl_file
        self.tree = ET.parse(wadl_file)
        self.root = self.tree.getroot()
        self.methods: Dict = {}
    
    def parse(self) -> Dict:
        resources = self.root.findall('.//wadl:resource', self.NAMESPACES)
        
        for resource in resources:
            resource_id = resource.get('id')
            sample_path = resource.get('{%s}samplePath' % self.NAMESPACES['wx'])
            
            methods = resource.findall('wadl:method', self.NAMESPACES)
            for method in methods:
                method_id = method.get('id')
                method_name = method.get('name')
                
                spec = self._extract_method_spec(method, resource_id, sample_path)
                spec['resource_id'] = resource_id
                spec['sample_path'] = sample_path
                
                self.methods[method_id] = spec
        
        return self.methods
    
    def _extract_method_spec(self, method_elem, resource_id: str, sample_path: str) -> Dict:
        method_name = method_elem.get('name')
        responses = self._extract_responses(method_elem)
        response_types = self._extract_response_types(method_elem)
        
        return {
            'method': method_name,
            'responses': responses,
            'response_types': response_types,
        }
    
    def _extract_responses(self, method_elem) -> List[Dict]:
        responses = []
        response_elems = method_elem.findall('wadl:response', self.NAMESPACES)
        
        if not response_elems:
            responses.append({'status': 200, 'headers': [], 'media_types': []})
        else:
            for resp in response_elems:
                status = resp.get('status', '200')
                headers = self._extract_headers(resp)
                media_types = self._extract_media_types(resp)
                responses.append({'status': int(status), 'headers': headers, 'media_types': media_types})
        
        return responses
    
    def _extract_headers(self, response_elem) -> List[str]:
        headers = []
        params = response_elem.findall('wadl:param', self.NAMESPACES)
        
        for param in params:
            style = param.get('style')
            if style == 'header':
                name = param.get('name')
                required = param.get('required', 'false').lower() == 'true'
                if name and required:
                    headers.append(name)
        
        return headers
    
    def _extract_media_types(self, response_elem) -> List[str]:
        media_types = []
        reps = response_elem.findall('wadl:representation', self.NAMESPACES)
        
        for rep in reps:
            media_type = rep.get('mediaType')
            if media_type:
                media_types.append(media_type)
        
        return media_types
    
    def _extract_response_types(self, method_elem) -> List[str]:
        types = []
        responses = method_elem.findall('wadl:response', self.NAMESPACES)
        
        for resp in responses:
            reps = resp.findall('wadl:representation', self.NAMESPACES)
            for rep in reps:
                element = rep.get('element')
                if element and element.startswith('sep:'):
                    element_type = element.replace('sep:', '')
                    if element_type not in types:
                        types.append(element_type)
        
        return types


def count_methods_in_wadl(wadl_dir: str) -> int:
    count = 0
    for wadl_file in Path(wadl_dir).glob('*.wadl'):
        parser = WADLParser(str(wadl_file))
        parser.parse()
        count += len(parser.methods)
    return count


if __name__ == '__main__':
    wadl_dir = '/home/tylor/dev/egot/wadl'
    print(f"Total methods in all WADLs: {count_methods_in_wadl(wadl_dir)}")
