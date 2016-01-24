package com.bullercodeworks.magopie.adapters;

import android.content.Context;
import android.util.Log;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;
import android.widget.Toast;

import com.bullercodeworks.magopie.R;

import java.util.List;

import go.magopie.Magopie;

public class ResultAdapter extends ArrayAdapter<Magopie.Torrent> {
  private Context ctx;
  public ResultAdapter(List<Magopie.Torrent> tList, Context _ctx) {
    super(_ctx, R.layout.torrent_list_row, R.id.torrentTitle, tList);
    ctx = _ctx;
  }

  public class ViewHolder {
    public TextView torrentTitle;
    public TextView torrentFile;

    public Magopie.Torrent torrent;

    ViewHolder(View row) {
      this.torrentFile = (TextView)row.findViewById(R.id.torrentFile);
      this.torrentTitle = (TextView)row.findViewById(R.id.torrentTitle);
    }
  }

  @Override
  public View getView(int position, View convertView, ViewGroup parent) {
    View row = super.getView(position, convertView, parent);
    ViewHolder holder = (ViewHolder)row.getTag();
    if(holder == null) {
      holder = new ViewHolder(row);
      row.setTag(holder);
    }
    final Magopie.Torrent wrk = getItem(position);
    if(wrk != null) {
      holder.torrent = wrk;
      holder.torrentFile.setText(wrk.getFileURL());
      holder.torrentTitle.setText(wrk.getTitle());
    }

    row.setOnClickListener(new View.OnClickListener() {
      @Override
      public void onClick(View view) {
        if(Magopie.Download(((ViewHolder)view.getTag()).torrent)) {
          Toast.makeText(ctx, "Torrent download triggered", Toast.LENGTH_SHORT).show();
        } else {
          Toast.makeText(ctx, "Failed to trigger download!", Toast.LENGTH_SHORT).show();
        }
      }
    });
    return row;
  }
}