"use strict";
(function(CodeMirror) {
    // Modified from the sample simplemode
    CodeMirror.defineSimpleMode("goasm", {
        start: [
            {regex: /"(?:[^\\]|\\.)*?(?:"|$)/, token: "string"},
            {regex: /(?:getc|out|puts|in|halt|add|and|nop|br|brn|brz|brp|brnp|brnz|brzp|brnzp|jmp|jmpr|jsr|jsrr|ld|ldi|ldr|lea|not|ret|rti|st|sti|str|trap)\b/,
             token: "keyword"}, // instructions
            {regex: /\.[^\s]+/, token: "atom"}, // assembly directive
            {regex: /(?:\$|x|X)[a-f\d]+|#?\d+|(?:%|b|B)[01]+/,
             token: "number"},
            {regex: /;.*/, token: "comment"},
            {regex: /(?:r\d|pc|ir)\b/, token: "variable-2"}, // registers
        ]
    });

    CodeMirror.defineMIME("text/x-goasm", "goasm");
    CodeMirror.defineMIME("text/x-goasm", { name: "goasm" });
})(CodeMirror);
const goasm = (function() {
    var gocm;
    function resizeToFit(elem) {
        elem.style.height = elem.scrollHeight + 'px';
    }
    function displayAsm(asm) {
        const out = document.getElementById('asm');
        out.innerHTML = asm;
        resizeToFit(out);
    }
    function getAsm(evt) {
        evt.preventDefault();
        const url = evt.target.action;
        const method = evt.target.method.toUpperCase();
        const fd = new FormData(evt.target);
        const req = new XMLHttpRequest;
        const async = true;
        req.onload = function(evt) {
            if(req.status === 200) {
                displayAsm(req.responseText);
            }
        }
        req.open(method, url, async);
        req.send(fd);
    }
    function init() {
        const gota = document.getElementById('gocode');
        gocm = CodeMirror.fromTextArea(gota, {
            theme: 'default left',
            lineNumbers: true,
            mode: 'text/x-go',
        });
        const form = document.forms['asm-form'];
        const useCapture = true;
        form.addEventListener('submit', getAsm, useCapture);
    }
    return {init:init};
})();
goasm.init();
