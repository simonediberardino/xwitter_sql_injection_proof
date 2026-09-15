# -*- coding: utf-8 -*-
"""Reformat section 6 (last point) in the homework Word document."""
from copy import deepcopy

from docx import Document
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.oxml import OxmlElement
from docx.oxml.ns import qn
from docx.shared import Cm, Pt, RGBColor

from pathlib import Path

DOC = Path(__file__).resolve().parent / "Homework_Sicurezza_Simone_Di_Berardino_2085977.docx"
SHADE = "F3F3F3"
BLACK = RGBColor(0x1A, 0x1A, 0x1A)


def set_run_font(run, name="Calibri", size=12, bold=False, italic=False):
    run.font.name = name
    run.font.size = Pt(size)
    run.font.bold = bold
    run.font.italic = italic
    run.font.color.rgb = BLACK
    rPr = run._element.get_or_add_rPr()
    rFonts = rPr.find(qn("w:rFonts"))
    if rFonts is None:
        rFonts = OxmlElement("w:rFonts")
        rPr.append(rFonts)
    rFonts.set(qn("w:ascii"), name)
    rFonts.set(qn("w:hAnsi"), name)


def shade_cell(cell, fill=SHADE):
    tcPr = cell._tc.get_or_add_tcPr()
    shd = OxmlElement("w:shd")
    shd.set(qn("w:fill"), fill)
    shd.set(qn("w:val"), "clear")
    tcPr.append(shd)


def set_cell_border(cell):
    tcPr = cell._tc.get_or_add_tcPr()
    tcBorders = OxmlElement("w:tcBorders")
    for edge in ("top", "left", "bottom", "right"):
        element = OxmlElement(f"w:{edge}")
        element.set(qn("w:val"), "single")
        element.set(qn("w:sz"), "4")
        element.set(qn("w:color"), "DDDDDD")
        tcBorders.append(element)
    tcPr.append(tcBorders)


def clear_paragraph(paragraph):
    paragraph.clear()


def write_body(paragraph, text, justify=True):
    try:
        paragraph.style = "Normal"
    except KeyError:
        pass
    clear_paragraph(paragraph)
    paragraph.alignment = WD_ALIGN_PARAGRAPH.JUSTIFY if justify else WD_ALIGN_PARAGRAPH.LEFT
    paragraph.paragraph_format.space_after = Pt(8)
    paragraph.paragraph_format.line_spacing = 1.15
    run = paragraph.add_run(text)
    set_run_font(run, "Calibri", 12)


def heading_style(doc, level):
    target = f"Heading {level}"
    for p in doc.paragraphs:
        if p.style and p.style.name == target:
            return p.style
    return doc.paragraphs[0].style


def write_heading1(doc, paragraph, text):
    paragraph.style = heading_style(doc, 1)
    clear_paragraph(paragraph)
    run = paragraph.add_run(text)
    set_run_font(run, "Arial", 16, bold=True)


def write_heading2(doc, paragraph, text):
    paragraph.style = heading_style(doc, 2)
    clear_paragraph(paragraph)
    run = paragraph.add_run(text)
    set_run_font(run, "Arial", 13, bold=True)


def add_code_table_after(paragraph, code_text):
    """Insert a one-cell code table immediately after paragraph."""
    table = paragraph._parent.add_table(rows=1, cols=1, width=Cm(16))
    cell = table.cell(0, 0)
    shade_cell(cell)
    set_cell_border(cell)
    cell.text = ""
    p = cell.paragraphs[0]
    p.paragraph_format.space_before = Pt(0)
    p.paragraph_format.space_after = Pt(0)
    run = p.add_run(code_text.strip("\n"))
    set_run_font(run, "Consolas", 9)

    # Move table element to sit right after the paragraph
    p_el = paragraph._p
    tbl_el = table._tbl
    p_el.addnext(tbl_el)
    return table


def insert_paragraph_after(paragraph, text="", style_name="Normal"):
    new_p = OxmlElement("w:p")
    paragraph._p.addnext(new_p)
    from docx.text.paragraph import Paragraph

    new_para = Paragraph(new_p, paragraph._parent)
    if style_name:
        new_para.style = style_name
    if text:
        run = new_para.add_run(text)
        set_run_font(run, "Calibri", 12)
    return new_para


def main():
    doc = Document(DOC)

    # Section 6 starts at paragraph index 40 in the current file.
    idx = None
    for i, p in enumerate(doc.paragraphs):
        if p.text.strip().lower().startswith("6.") and "modifica" in p.text.lower():
            idx = i
            break
    if idx is None:
        raise SystemExit("Section 6 heading not found")

    p_title = doc.paragraphs[idx]
    write_heading1(doc, p_title, "6. Modifica del contenuto del database")

    # Reuse following paragraphs; create extras if needed
    while len(doc.paragraphs) <= idx + 6:
        insert_paragraph_after(doc.paragraphs[-1])

    paragraphs = doc.paragraphs
    p_intro = paragraphs[idx + 1]
    p_mecc = paragraphs[idx + 2]
    p_payload_intro = paragraphs[idx + 3]
    p_after_code = paragraphs[idx + 4]

    write_body(
        p_intro,
        "Oltre alla Search user, anche la pagina Profile consente di alterare i dati "
        "memorizzati nel database. In laboratorio la modifica della bio avviene tramite "
        "POST /api/auth/profile: il valore del campo bio viene concatenato in una UPDATE "
        "SQL separata, senza parametri. Non è richiesta alcuna password per cambiare "
        "username, display name o bio, quindi chi controlla il payload inviato al server "
        "può sfruttare lo stesso punto di ingresso per un attacco piggybacked.",
    )

    write_heading2(doc, p_mecc, "6.1 Meccanismo dell'injection")

    write_body(
        p_payload_intro,
        "La query vulnerabile ha la forma UPDATE users SET bio = '<input>' WHERE id = <id>. "
        "Chiudendo la stringa con un apice e aggiungendo una seconda assegnazione, seguita "
        "da un commento di fine riga (--), il motore SQL esegue anche l'aggiornamento di "
        "password_hash. L'aggiornamento colpisce la riga dell'utente autenticato (clausola "
        "WHERE legata all'id), ma il payload dimostra che un singolo campo testuale non "
        "sanitizzato basta a riscrivere colonne sensibili.",
    )

    # Clear old payload line and insert structured content after p_payload_intro
    if len(paragraphs) > idx + 5:
        old_payload = paragraphs[idx + 5]
        write_heading2(doc, old_payload, "6.2 Payload usato in Profile")
    else:
        old_payload = insert_paragraph_after(p_after_code, style_name="Heading 2")
        write_heading2(doc, old_payload, "6.2 Payload usato in Profile")

    p_explain = insert_paragraph_after(old_payload)
    write_body(
        p_explain,
        "Nel campo Bio della pagina Profile è stato inserito il seguente input. "
        "La parte dopo l'apice iniziale viene interpretata come SQL aggiuntivo; "
        "i due trattini commentano il resto della query generata dal backend.",
    )

    code_table = add_code_table_after(
        p_explain,
        "you have been hacked', password_hash = 'password_hacked'--",
    )

    p_result = insert_paragraph_after(p_explain)
    # Table was inserted after p_explain; move result paragraph after table
    tbl_el = code_table._tbl
    p_result._p.getparent().remove(p_result._p)
    tbl_el.addnext(p_result._p)

    write_body(
        p_result,
        "Effetto osservato: la bio dell'account diventa you have been hacked e il campo "
        "password_hash viene impostato al valore letterale password_hacked (non un hash "
        "bcrypt valido). Gli account interessati non possono più autenticarsi con la "
        "password precedente. Questo passo soddisfa l'obiettivo «modificare il contenuto "
        "del database» e illustra la violazione dell'integrità (Integrity) dei dati.",
    )

    p_cia = insert_paragraph_after(p_result)
    write_body(
        p_cia,
        "In sintesi, sulla stessa applicazione Xwitter: la Search user ha permesso "
        "enumerazione ed esfiltrazione (Confidentiality); la Profile ha permesso "
        "alterazione di bio e password_hash (Integrity). Entrambi gli scenari restano "
        "simulazioni didattiche su codice volutamente vulnerabile.",
    )

    # Remove leftover empty or duplicate paragraphs at old indices if any still contain raw payload only
    for p in doc.paragraphs:
        t = p.text.strip()
        if t == "you have been hacked', password_hash = 'password_hacked'--" and p._p != p_explain._p:
            clear_paragraph(p)

    doc.save(DOC)
    print("Updated:", DOC)


if __name__ == "__main__":
    main()
