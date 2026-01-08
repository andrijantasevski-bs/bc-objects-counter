enum 50100 "Sample Status"
{
    Extensible = true;

    value(0; Open) { Caption = 'Open'; }
    value(1; Closed) { Caption = 'Closed'; }
}

report 50100 "Sample Report"
{
    DefaultLayout = RDLC;

    dataset
    {
        dataitem("Sample Customer"; "Sample Customer")
        {
        }
    }
}

tableextension 50100 "Customer Ext" extends Customer
{
    fields
    {
        field(50100; "Custom Field"; Text[50])
        {
            Caption = 'Custom Field';
        }
    }
}
